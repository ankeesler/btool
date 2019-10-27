#include "app/collector/cc/inc.h"

#include <string>

#include "app/collector/cc/properties.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "core/log.h"
#include "node/node.h"
#include "util/string/string.h"

namespace btool::app::collector::cc {

static void HandleInclude(::btool::app::collector::Store *s,
                          ::btool::node::Node *n, const std::string &include);

void Inc::OnSet(::btool::app::collector::Store *s, const std::string &name) {
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

  ip_->ParseIncludes(
      [s, n](const std::string &include) { HandleInclude(s, n, include); });
}

static void HandleInclude(::btool::app::collector::Store *s,
                          ::btool::node::Node *n, const std::string &include) {
  ::btool::node::Node *d = nullptr;
  std::string include_path;

  for (auto it : *s) {
    auto sn = it.second;
    DEBUG("does node %s end in include %s\n", sn->Name().c_str(),
          include.c_str());
    if (::btool::util::string::HasSuffix(sn->Name().c_str(), include.c_str())) {
      d = sn;

      std::size_t index = sn->Name().rfind(include);
      if (index == 0) {
        include_path = ".";
      } else {
        include_path = sn->Name().substr(0, index);
      }
      DEBUG("yes, and the include path is %s\n", include_path.c_str());
      break;
    }
  }

  if (d == nullptr) {
    throw "ahhhh noooo exceptions are bad!!!";
  }

  n->dependencies()->push_back(d);
  Properties::AddIncludePath(n->property_store(), include_path);
}

};  // namespace btool::app::collector::cc
