#include "app/collector/cc/inc.h"

#include <cassert>

#include <algorithm>
#include <string>

#include "app/collector/cc/properties.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "core/log.h"
#include "node/node.h"
#include "util/string/string.h"
#include "util/util.h"

namespace btool::app::collector::cc {

static bool HandleInclude(::btool::app::collector::Store *s,
                          ::btool::node::Node *n, const std::string &include);

void Inc::OnNotify(::btool::app::collector::Store *s, const std::string &name) {
  auto n = s->Get(name);
  if (n == nullptr) {
    return;
  }

  if (!::btool::app::collector::Properties::Local(n->property_store())) {
    return;
  }

  bool c = ::btool::util::string::HasSuffix(name.c_str(), ".c");
  bool cc = ::btool::util::string::HasSuffix(name.c_str(), ".cc");
  bool h = ::btool::util::string::HasSuffix(name.c_str(), ".h");
  if (!c && !cc && !h) {
    return;
  }

  bool updated = false;
  auto err =
      ip_->ParseIncludes(name, [s, n, &updated](const std::string &include) {
        bool new_stuff = HandleInclude(s, n, include);
        updated = updated || new_stuff;
      });
  if (err) {
    DEBUG("failed to parse includes %s\n", err.Msg());
    assert(0);
  }

  if (updated) {
    Notify(s, name);
  }
}

static bool HandleInclude(::btool::app::collector::Store *s,
                          ::btool::node::Node *n, const std::string &include) {
  DEBUG("handling include %s\n", include.c_str());

  ::btool::node::Node *d = nullptr;
  std::string include_path;

  for (auto it : *s) {
    auto sn = it.second;
    DEBUG("does node %s end in include %s\n", sn->name().c_str(),
          include.c_str());
    if (::btool::util::string::HasSuffix(sn->name().c_str(), include.c_str())) {
      d = sn;

      std::size_t index = sn->name().rfind(include);
      if (index == 0) {
        include_path = ".";
      } else {
        include_path = sn->name().substr(0, index);
      }
      DEBUG("yes, and the include path is %s\n", include_path.c_str());
      break;
    }
  }

  if (d == nullptr) {
    DEBUG("cannot resolve include %s for node %s\n", include.c_str(),
          n->name().c_str());
    return false;
  }

  if (::btool::util::Contains(*n->dependencies(), d)) {
    return false;
  }

  n->dependencies()->push_back(d);

  auto include_paths = Properties::IncludePaths(n->property_store());
  bool needs_add = true;
  if (include_paths != nullptr) {
    auto it =
        std::find(include_paths->begin(), include_paths->end(), include_path);
    if (it != include_paths->end()) {
      needs_add = false;
    }
  }

  if (needs_add) {
    Properties::AddIncludePath(n->property_store(), include_path);
  }

  return true;
}

};  // namespace btool::app::collector::cc
