#include "fs_collectini.h"

#include <dirent.h>
#include <errno.h>
#include <cstring>

#include <functional>
#include <vector>

#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "err.h"
#include "log.h"
#include "node/node.h"
#include "util/fs/fs.h"
#include "util/string/string.h"

namespace btool::app::collector::fs {

void FSCollectini::Collect(::btool::app::collector::Store *s) {
  std::vector<::btool::node::Node *> nodes;

  ::btool::util::fs::Walk(root_, [&](const std::string &path) {
    auto is_dir = ::btool::util::fs::IsDir(path);

    if (!is_dir && (::btool::util::string::HasSuffix(path.c_str(), ".c") ||
                    ::btool::util::string::HasSuffix(path.c_str(), ".cc") ||
                    ::btool::util::string::HasSuffix(path.c_str(), ".h"))) {
      auto n = s->Put(path.c_str());
      ::btool::app::collector::Properties::SetLocal(n->property_store(), true);
      nodes.push_back(n);
    }
  });

  for (auto n : nodes) {
    Notify(s, n->name());
  }
}

};  // namespace btool::app::collector::fs
