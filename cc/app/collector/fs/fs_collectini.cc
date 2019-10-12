#include "fs_collectini.h"

#include <dirent.h>
#include <errno.h>
#include <cstring>

#include <functional>
#include <iostream>
#include <memory>

#include "core/err.h"
#include "core/log.h"
#include "node/store.h"

namespace btool::app::collector::fs {

static ::btool::core::VoidErr Walk(const std::string &root,
                                   std::function<void(const std::string &)> f);
static bool HasExt(const char *file, const char *ext);

::btool::core::VoidErr FSCollectini::Collect(::btool::node::Store *s) {
  return Walk(root_, [&](const std::string &path) { s->Put(path.c_str()); });
}

static ::btool::core::VoidErr Walk(const std::string &root,
                                   std::function<void(const std::string &)> f) {
  DEBUG("walk %s\n", root.c_str());

  std::vector<std::string> children;
  ::DIR *dir = ::opendir(root.c_str());
  if (dir == nullptr) {
    DEBUG("opendir: %s\n", ::strerror(errno));
    return ::btool::core::VoidErr::Failure("opendir");
  }

  struct dirent *dirent = nullptr;
  while ((dirent = readdir(dir)) != nullptr) {
    const char *d_name = dirent->d_name;
    if ((dirent->d_type & DT_DIR) != 0) {
      if (::strcmp(d_name, ".") != 0 && ::strcmp(d_name, "..") != 0) {
        children.push_back(d_name);
      }
    } else if ((dirent->d_type & DT_REG) != 0) {
      if (HasExt(d_name, ".c") || HasExt(d_name, ".cc") ||
          HasExt(d_name, ".h")) {
        std::string file(root + '/' + d_name);
        DEBUG("visit %s\n", file.c_str());
        f(file);
      }
    }
  }

  if (::closedir(dir) == -1) {
    DEBUG("closedir: %s\n", ::strerror(errno));
    return ::btool::core::VoidErr::Failure("closedir");
  }

  for (auto child : children) {
    std::string new_root(root + '/' + child);
    auto err = Walk(new_root, f);
    if (err) {
      return err;
    }
  }

  return ::btool::core::VoidErr::Success();
}

static bool HasExt(const char *file, const char *ext) {
  size_t file_len = ::strlen(file);
  size_t ext_len = ::strlen(ext);
  return (file_len > ext_len &&
          (::memcmp(file + file_len - ext_len, ext, ext_len) == 0));
}

};  // namespace btool::app::collector::fs
