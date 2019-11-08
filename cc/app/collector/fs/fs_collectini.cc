#include "fs_collectini.h"

#include <dirent.h>
#include <errno.h>
#include <cstring>

#include <functional>
#include <iostream>
#include <memory>

#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "core/err.h"
#include "core/log.h"
#include "util/fs/fs.h"
#include "util/string/string.h"

namespace btool::app::collector::fs {

void FSCollectini::Collect(::btool::app::collector::Store *s) {
  ::btool::util::fs::Walk(
      root_, [&](const std::string &path) -> ::btool::core::VoidErr {
        auto err = ::btool::util::fs::IsFile(path);
        if (err) {
          return ::btool::core::VoidErr::Failure(err.Msg());
        }

        if (err.Ret() &&
            (::btool::util::string::HasSuffix(path.c_str(), ".c") ||
             ::btool::util::string::HasSuffix(path.c_str(), ".cc") ||
             ::btool::util::string::HasSuffix(path.c_str(), ".h"))) {
          auto n = s->Put(path.c_str());
          ::btool::app::collector::Properties::SetLocal(n->property_store(),
                                                        true);
          Notify(s, n->name());
        }

        return ::btool::core::VoidErr::Success();
      });
}

};  // namespace btool::app::collector::fs
