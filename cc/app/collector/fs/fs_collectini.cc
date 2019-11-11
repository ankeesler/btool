#include "fs_collectini.h"

#include <dirent.h>
#include <errno.h>
#include <cstring>

#include <functional>
#include <vector>

#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "core/err.h"
#include "core/log.h"
#include "node/node.h"
#include "util/fs/fs.h"
#include "util/string/string.h"

namespace btool::app::collector::fs {

void FSCollectini::Collect(::btool::app::collector::Store *s) {
  std::vector<::btool::node::Node *> nodes;

  ::btool::util::fs::Walk(
      root_, [&](const std::string &path) -> ::btool::core::VoidErr {
        auto err = ::btool::util::fs::IsDir(path);
        if (err) {
          return ::btool::core::VoidErr::Failure(err.Msg());
        }

        if (!err.Ret() &&
            (::btool::util::string::HasSuffix(path.c_str(), ".c") ||
             ::btool::util::string::HasSuffix(path.c_str(), ".cc") ||
             ::btool::util::string::HasSuffix(path.c_str(), ".h"))) {
          auto n = s->Put(path.c_str());
          ::btool::app::collector::Properties::SetLocal(n->property_store(),
                                                        true);
          nodes.push_back(n);
        }

        return ::btool::core::VoidErr::Success();
      });

  for (auto n : nodes) {
    Notify(s, n->name());
  }
}

};  // namespace btool::app::collector::fs
