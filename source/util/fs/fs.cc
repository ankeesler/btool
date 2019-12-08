#include "fs.h"

#include <dirent.h>
#include <errno.h>
#include <sys/stat.h>
#include <unistd.h>
#include <cstdio>
#include <cstring>

#include <fstream>
#include <functional>
#include <list>
#include <stack>

#include "err.h"
#include "log.h"

namespace btool::util::fs {

const int kIFBufSizeLog = 10;  // 1KB

std::string Base(const std::string &path) {
  std::size_t found = path.rfind('/');
  if (found == std::string::npos) {
    // single-name-directory-or-file.c
    return path;
  } else if (found == 0) {
    // /some-root-directory-or-file.c
    return "/";
  } else if (found == path.size() - 1) {
    // path/to/directory-with-slash/
    found = path.rfind('/', found - 1);
    return path.substr(found + 1, path.size() - 1 - found - 1);
  } else {
    // path/to/directory-or-file-without-slash.c
    return path.substr(found + 1, path.size() - found - 1);
  }
}

std::string Dir(const std::string &path) {
  std::size_t found = path.rfind('/');
  if (found == std::string::npos) {
    // single-name-directory-or-file.c
    return ".";
  } else if (found == 0) {
    // /some-root-directory-or-file.c
    return "/";
  } else if (found == path.size() - 1) {
    // path/to/directory-with-slash/
    found = path.rfind('/', found - 1);
    return path.substr(0, found);
  } else {
    // path/to/directory-or-file-without-slash.c
    return path.substr(0, found);
  }
}

std::string Join(const std::string &one, const std::string &two) {
  return one + '/' + two;
}

std::string Ext(const std::string &path) {
  auto path_base = Base(path);
  std::size_t index = path_base.rfind('.');
  if (index == std::string::npos || index == 0) {
    return "";
  } else {
    return path_base.substr(index);
  }
  return "";
}

std::string TempDir() {
  char s[] = "/tmp/btool_XXXXXX";
  char *dir = ::mkdtemp(s);
  if (dir == NULL) {
    THROW_ERR("mkdtmp: " + std::string(::strerror(errno)));
  } else {
    return dir;
  }
}

std::string ReadFile(const std::string &path) {
  FILE *f = ::fopen(path.c_str(), "r");
  if (f == nullptr) {
    THROW_ERR("fopen " + path + ": " + std::string(::strerror(errno)));
  }

  std::string content;
  while (true) {
    const int buf_size = 1 << kIFBufSizeLog;
    char buf[buf_size];
    ::size_t count = ::fread(buf, 1, buf_size, f);
    content.append(buf, count);
    if (::ferror(f)) {
      THROW_ERR("fread: " + std::string(::strerror(errno)));
    } else if (::feof(f)) {
      break;
    }
  }

  ::fclose(f);

  return content;
}

void WriteFile(const std::string &path, const std::string &content) {
  std::ofstream ofs(path);
  if (!ofs) {
    THROW_ERR("open: " + std::string(::strerror(errno)));
  }

  ofs << content;
  if (!ofs) {
    THROW_ERR("write: " + std::string(::strerror(errno)));
  }
}

void RemoveAll(const std::string &path) {
  auto exists = Exists(path);
  if (!exists) {
    return;
  }

  auto is_dir = IsDir(path);
  if (!is_dir) {
    if (::remove(path.c_str()) == -1) {
      THROW_ERR("remove: " + std::string(::strerror(errno)));
    }
    return;
  }

  Walk(path, [](const std::string &path) {
    if (::remove(path.c_str()) == -1) {
      THROW_ERR("remove: " + std::string(::strerror(errno)));
    }
  });
}

void Mkdir(const std::string &path) {
  if (::mkdir(path.c_str(), 0700) == -1) {
    THROW_ERR("mkdir: " + std::string(::strerror(errno)));
  }
}

void MkdirAll(const std::string &path) {
  std::stack<std::string> dirs;
  dirs.push(path);
  while (dirs.top() != "." && dirs.top() != "/") {
    dirs.push(Dir(dirs.top()));
  }

  while (!dirs.empty()) {
    if (!Exists(dirs.top())) {
      Mkdir(dirs.top());
    }
    dirs.pop();
  }
}

bool Exists(const std::string &path) {
  bool exists = true;
  struct ::stat s;
  if (::stat(path.c_str(), &s) == -1) {
    if (errno == ENOENT) {
      exists = false;
    } else {
      THROW_ERR("Exists " + path + ": stat: " + std::string(::strerror(errno)));
    }
  }
  return exists;
}

bool IsDir(const std::string &path) {
  bool is_dir = false;
  struct ::stat s;
  if (::stat(path.c_str(), &s) == -1) {
    if (errno != ENOENT) {
      THROW_ERR("IsDir " + path + ": stat: " + std::string(::strerror(errno)));
    }
  } else {
    is_dir = (((s.st_mode & S_IFMT) & S_IFDIR) != 0);
  }
  return is_dir;
}

void Walk(const std::string &root, std::function<void(const std::string &)> f) {
  DEBUG("walk %s\n", root.c_str());

  ::DIR *dir = ::opendir(root.c_str());
  if (dir == nullptr) {
    THROW_ERR("opendir '" + root + "': " + std::string(::strerror(errno)));
  }

  std::list<std::string> dir_children, file_children;
  struct dirent *dirent = nullptr;
  while ((dirent = readdir(dir)) != nullptr) {
    const char *d_name = dirent->d_name;
    if ((dirent->d_type & DT_DIR) != 0) {
      if (::strcmp(d_name, ".") != 0 && ::strcmp(d_name, "..") != 0) {
        dir_children.push_back(d_name);
      }
    } else if ((dirent->d_type & DT_REG) != 0) {
      file_children.push_back(d_name);
    }
  }

  if (::closedir(dir) == -1) {
    THROW_ERR("closedir: " + std::string(::strerror(errno)));
  }

  dir_children.sort();
  for (const auto &child : dir_children) {
    Walk(Join(root, child), f);
  }

  file_children.sort();
  for (const auto &child : file_children) {
    f(Join(root, child));
  }

  f(root);
}

}  // namespace btool::util::fs
