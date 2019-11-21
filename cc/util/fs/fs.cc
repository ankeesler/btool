#include "util/fs/fs.h"

#include <dirent.h>
#include <errno.h>
#include <sys/stat.h>
#include <unistd.h>
#include <cstdio>
#include <cstring>

#include <fstream>
#include <functional>
#include <list>

#include "log.h"

namespace btool::util::fs {

const int kIFBufSizeLog = 10;  // 1KB

std::string Base(const std::string &path) {
  std::size_t found = path.rfind('/');
  if (found == std::string::npos) {
    return path;
  } else {
    return std::string(path.c_str() + found + 1, path.size() - found - 1);
  }
}

std::string Dir(const std::string &path) {
  std::size_t found = path.rfind('/');
  if (found == std::string::npos) {
    return path;
  } else {
    return std::string(path.c_str(), found);
  }
}

std::string Join(const std::string &one, const std::string &two) {
  return one + '/' + two;
}

std::string Ext(const std::string &path) {
  std::size_t index = path.rfind('.');
  if (index == std::string::npos || index == 0) {
    return "";
  } else {
    return path.substr(index);
  }
  return "";
}

::btool::Err<std::string> TempDir() {
  char s[] = "/tmp/btool_XXXXXX";
  char *dir = ::mkdtemp(s);
  if (dir == NULL) {
    return ::btool::Err<std::string>::Failure("mkdtemp");
  } else {
    return ::btool::Err<std::string>::Success(dir);
  }
}

::btool::Err<std::string> ReadFile(const std::string &path) {
  FILE *f = ::fopen(path.c_str(), "r");
  if (f == nullptr) {
    return ::btool::Err<std::string>::Failure(strerror(errno));
  }

  std::string content;
  while (true) {
    const int buf_size = 1 << kIFBufSizeLog;
    char buf[buf_size];
    ::size_t count = ::fread(buf, 1, buf_size, f);
    content.append(buf, count);
    if (ferror(f)) {
      return ::btool::Err<std::string>::Failure(strerror(errno));
    } else if (feof(f)) {
      break;
    }
  }

  fclose(f);

  return ::btool::Err<std::string>::Success(content);
}

::btool::VoidErr WriteFile(const std::string &path,
                           const std::string &content) {
  std::ofstream ofs(path);
  if (!ofs) {
    return ::btool::VoidErr::Failure(strerror(errno));
  }

  ofs << content;
  if (!ofs) {
    return ::btool::VoidErr::Failure(strerror(errno));
  }

  return ::btool::VoidErr::Success();
}

::btool::VoidErr RemoveAll(const std::string &path) {
  auto err = Exists(path);
  if (err) {
    return ::btool::VoidErr::Failure(err.Msg());
  } else if (!err.Ret()) {
    return ::btool::VoidErr::Success();
  }

  err = IsDir(path);
  if (err) {
    return ::btool::VoidErr::Failure(err.Msg());
  } else if (!err.Ret()) {
    if (::remove(path.c_str()) == -1) {
      return ::btool::VoidErr::Failure(strerror(errno));
    }
    return ::btool::VoidErr::Success();
  }

  return Walk(path, [](const std::string &path) -> ::btool::VoidErr {
    if (::remove(path.c_str()) == -1) {
      return ::btool::VoidErr::Failure(strerror(errno));
    }
    return ::btool::VoidErr::Success();
  });
}

::btool::VoidErr Mkdir(const std::string &path) {
  if (::mkdir(path.c_str(), 0700) == -1) {
    return ::btool::VoidErr::Failure(strerror(errno));
  }
  return ::btool::VoidErr::Success();
}

::btool::Err<bool> Exists(const std::string &path) {
  bool exists = true;
  struct ::stat s;
  if (::stat(path.c_str(), &s) == -1) {
    if (errno == ENOENT) {
      exists = false;
    } else {
      return ::btool::Err<bool>::Failure(strerror(errno));
    }
  }
  return ::btool::Err<bool>::Success(exists);
}

::btool::Err<bool> IsDir(const std::string &path) {
  bool is_dir = false;
  struct ::stat s;
  if (::stat(path.c_str(), &s) == -1) {
    if (errno != ENOENT) {
      return ::btool::Err<bool>::Failure(strerror(errno));
    }
  } else {
    is_dir = (((s.st_mode & S_IFMT) & S_IFDIR) != 0);
  }
  return ::btool::Err<bool>::Success(is_dir);
}

::btool::VoidErr Walk(const std::string &root,
                      std::function<::btool::VoidErr(const std::string &)> f) {
  DEBUG("walk %s\n", root.c_str());

  ::DIR *dir = ::opendir(root.c_str());
  if (dir == nullptr) {
    DEBUG("opendir: %s\n", ::strerror(errno));
    return ::btool::VoidErr::Failure("opendir");
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
    DEBUG("closedir: %s\n", ::strerror(errno));
    return ::btool::VoidErr::Failure("closedir");
  }

  dir_children.sort();
  for (const auto &child : dir_children) {
    auto err = Walk(Join(root, child), f);
    if (err) {
      return err;
    }
  }

  file_children.sort();
  for (const auto &child : file_children) {
    auto err = f(Join(root, child));
    if (err) {
      return err;
    }
  }

  auto err = f(root);
  if (err) {
    return err;
  }

  return ::btool::VoidErr::Success();
}
};  // namespace btool::util::fs
