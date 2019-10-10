#include "fs_collectini.h"

#include <dirent.h>
#include <errno.h>
#include <sys/stat.h>
#include <cstring>

#include <functional>
#include <iostream>
#include <memory>

#include "core/err.h"
#include "core/log.h"
#include "node/store.h"

namespace btool::app::collector::fs {

static ::btool::core::VoidErr Walk(std::string *root,
                                   std::function<void(const std::string &)> f);

::btool::core::VoidErr FSCollectini::Collect(::btool::node::Store *s) {
  return Walk(new std::string(root_),
              [&](const std::string &path) { s->Put(path.c_str()); });
}

static ::btool::core::VoidErr Walk(std::string *root,
                                   std::function<void(const std::string &)> f) {
  DEBUG("walk %s\n", root->c_str());

  struct ::stat s;
  if (::stat(root->c_str(), &s) == -1) {
    DEBUG("lstat %s: %s\n", root->c_str(), ::strerror(errno));
    return ::btool::core::VoidErr::Failure("couldn't lstat node");
  }

  std::vector<const char *> children;
  if ((s.st_mode & S_IFDIR) != 0) {
    ::DIR *dir = ::opendir(root->c_str());
    if (dir == nullptr) {
      DEBUG("opendir: %s\n", ::strerror(errno));
      return ::btool::core::VoidErr::Failure("opendir");
    }

    struct dirent *dirent = nullptr;
    while ((dirent = readdir(dir)) != nullptr) {
      char *d_name = dirent->d_name;
      if (::strcmp(d_name, ".") != 0 && ::strcmp(d_name, "..") != 0) {
        children.push_back(d_name);
      }
    }

    if (::closedir(dir) == -1) {
      DEBUG("closedir: %s\n", ::strerror(errno));
      return ::btool::core::VoidErr::Failure("closedir");
    }
  } else {
    f(*root);
  }

  for (auto child : children) {
    auto new_root = new std::string(*root + "/" + child);
    auto err = Walk(new_root, f);
    if (err) {
      return err;
    }
  }

  delete root;

  return ::btool::core::VoidErr::Success();
}

};  // namespace btool::app::collector::fs
