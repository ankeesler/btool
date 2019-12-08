#include "app/collector/cc/inc.h"

#include <algorithm>
#include <string>

#include "app/collector/cc/properties.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "err.h"
#include "log.h"
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
  ip_->ParseIncludes(name, [s, n, &updated](const std::string &include) {
    bool new_stuff = HandleInclude(s, n, include);
    updated = updated || new_stuff;
  });

  if (updated) {
    Notify(s, name);
  }
}

static bool HandleInclude(::btool::app::collector::Store *s,
                          ::btool::node::Node *n, const std::string &include) {
  DEBUGS() << "handling include " << include << std::endl;

  ::btool::node::Node *d = nullptr;
  std::string include_path;

  for (auto it : *s) {
    auto sn = it.second;
    if (::btool::util::string::HasSuffix(sn->name().c_str(), include.c_str())) {
      d = sn;

      std::size_t index = sn->name().rfind(include);
      if (index == 0) {
        include_path = ".";
      } else {
        include_path = sn->name().substr(0, index);
      }
      DEBUGS() << "resolved include " << include << " to dependency "
               << d->name() << " and include path " << include_path
               << std::endl;
      break;
    }
  }

  if (d == nullptr) {
    DEBUGS() << "cannot resolve include " << include << "for node " << n->name()
             << std::endl;
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
