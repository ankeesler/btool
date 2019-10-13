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
#include "util/fs/fs.h"

namespace btool::app::collector::fs {

static bool HasExt(const char *file, const char *ext);

::btool::core::VoidErr FSCollectini::Collect(::btool::node::Store *s) {
  return ::btool::util::fs::Walk(
      root_, [&](const std::string &path) -> ::btool::core::VoidErr {
        auto err = ::btool::util::fs::IsFile(path);
        if (err) {
          return ::btool::core::VoidErr::Failure(err.Msg());
        }

        if (err.Ret() &&
            (HasExt(path.c_str(), ".c") || HasExt(path.c_str(), ".cc") ||
             HasExt(path.c_str(), ".h"))) {
          s->Put(path.c_str());
        }

        return ::btool::core::VoidErr::Success();
      });
}

static bool HasExt(const char *file, const char *ext) {
  size_t file_len = ::strlen(file);
  size_t ext_len = ::strlen(ext);
  return (file_len > ext_len &&
          (::memcmp(file + file_len - ext_len, ext, ext_len) == 0));
}

};  // namespace btool::app::collector::fs
