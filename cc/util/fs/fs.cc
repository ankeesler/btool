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

bool TempDir(std::string *ret_dir, std::string *ret_err) {
  char s[] = "/tmp/btool_XXXXXX";
  char *dir = ::mkdtemp(s);
  if (dir == NULL) {
    *ret_err = ::btool::WrapErr(::strerror(errno), "mkdtemp");
    return false;
  } else {
    *ret_dir = std::string(dir);
    return true;
  }
}

bool ReadFile(const std::string &path, std::string *ret_content,
              std::string *ret_err) {
  FILE *f = ::fopen(path.c_str(), "r");
  if (f == nullptr) {
    *ret_err = ::btool::WrapErr(::strerror(errno), "fopen");
    return false;
  }

  while (true) {
    const int buf_size = 1 << kIFBufSizeLog;
    char buf[buf_size];
    ::size_t count = ::fread(buf, 1, buf_size, f);
    ret_content->append(buf, count);
    if (::ferror(f)) {
      *ret_err = ::btool::WrapErr(::strerror(errno), "fread");
      return false;
    } else if (::feof(f)) {
      break;
    }
  }

  ::fclose(f);

  return true;
}

bool WriteFile(const std::string &path, const std::string &content,
               std::string *ret_err) {
  std::ofstream ofs(path);
  if (!ofs) {
    *ret_err = ::btool::WrapErr(::strerror(errno), "open");
    return false;
  }

  ofs << content;
  if (!ofs) {
    *ret_err = ::btool::WrapErr(::strerror(errno), "write");
    return false;
  }

  return true;
}

bool RemoveAll(const std::string &path, std::string *ret_err) {
  bool exists;
  std::string err;
  if (!Exists(path, &exists, &err)) {
    *ret_err = ::btool::WrapErr(err, "exists");
    return false;
  } else if (!exists) {
    return true;
  }

  bool is_dir;
  if (!IsDir(path, &is_dir, &err)) {
    *ret_err = ::btool::WrapErr(err, "is dir");
    return false;
  } else if (!is_dir) {
    if (::remove(path.c_str()) == -1) {
      *ret_err = ::btool::WrapErr(::strerror(errno), "remove");
      return false;
    }
    return true;
  }

  return Walk(
      path,
      [](const std::string &path, std::string *ret_err) -> bool {
        if (::remove(path.c_str()) == -1) {
          *ret_err = ::btool::WrapErr(::strerror(errno), "remove");
          return false;
        }
        return true;
      },
      ret_err);
}

bool Mkdir(const std::string &path, std::string *ret_err) {
  if (::mkdir(path.c_str(), 0700) == -1) {
    *ret_err = ::btool::WrapErr(::strerror(errno), "mkdir");
    return false;
  }
  return true;
}

bool Exists(const std::string &path, bool *ret_exists, std::string *ret_err) {
  *ret_exists = true;

  struct ::stat s;
  if (::stat(path.c_str(), &s) == -1) {
    if (errno == ENOENT) {
      *ret_exists = false;
    } else {
      *ret_err = ::btool::WrapErr(::strerror(errno), "stat");
      return false;
    }
  }

  return true;
}

bool IsDir(const std::string &path, bool *is, std::string *ret_err) {
  *is = false;

  struct ::stat s;
  if (::stat(path.c_str(), &s) == -1) {
    if (errno != ENOENT) {
      *ret_err = ::btool::WrapErr(::strerror(errno), "stat");
      return false;
    }
  } else {
    *is = (((s.st_mode & S_IFMT) & S_IFDIR) != 0);
  }

  return true;
}

bool Walk(const std::string &root,
          std::function<bool(const std::string &, std::string *)> f,
          std::string *ret_err) {
  DEBUG("walk %s\n", root.c_str());

  ::DIR *dir = ::opendir(root.c_str());
  if (dir == nullptr) {
    *ret_err = ::btool::WrapErr(::strerror(errno), "opendir");
    return false;
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
    *ret_err = ::btool::WrapErr(::strerror(errno), "opendir");
    return false;
  }

  dir_children.sort();
  for (const auto &child : dir_children) {
    if (!Walk(Join(root, child), f, ret_err)) {
      return false;
    }
  }

  std::string err;
  file_children.sort();
  for (const auto &child : file_children) {
    if (!f(Join(root, child), &err)) {
      *ret_err = ::btool::WrapErr(err, "handle");
      return false;
    }
  }

  if (!f(root, &err)) {
    *ret_err = ::btool::WrapErr(err, "handle");
    return false;
  }

  return true;
}
};  // namespace btool::util::fs
