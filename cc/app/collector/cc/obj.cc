#include "obj.h"

#include <string>
#include <vector>

#include "app/collector/cc/properties.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "core/log.h"
#include "node/node.h"
#include "util/string/string.h"

namespace btool::app::collector::cc {

void CollectIncludePaths(::btool::node::Node *n,
                         std::vector<std::string> *include_paths);

void Obj::OnNotify(::btool::app::collector::Store *s, const std::string &name) {
  auto d = s->Get(name);
  if (d == nullptr) {
    DEBUG("obj: drop %s (unknown name)\n", name.c_str());
    return;
  }

  if (!::btool::app::collector::Properties::Local(d->property_store())) {
    DEBUG("obj: drop %s (not local)\n", name.c_str());
    return;
  }

  bool c = ::btool::util::string::HasSuffix(name.c_str(), ".c");
  bool cc = ::btool::util::string::HasSuffix(name.c_str(), ".cc");
  if (!c && !cc) {
    DEBUG("obj: drop %s (not .c/.cc)\n", name.c_str());
    return;
  }

  std::string obj_name = name;
  if (c) {
    obj_name[obj_name.size() - 1] = 'o';
  } else {  // cc
    obj_name[obj_name.size() - 2] = 'o';
    obj_name.erase(obj_name.size() - 1, 1);
  }
  ::btool::node::Node *n = s->Put(obj_name);

  std::vector<std::string> include_paths;
  CollectIncludePaths(d, &include_paths);
  std::vector<std::string> flags;
  class ::btool::node::Node::Resolver *r;
  if (c) {
    r = rf_->NewCompileC(include_paths, flags);
  } else {  // cc
    r = rf_->NewCompileCC(include_paths, flags);
  }
  n->set_resolver(r);
  n->dependencies()->push_back(d);
  Notify(s, n->name());
}

void CollectIncludePaths(::btool::node::Node *n,
                         std::vector<std::string> *include_paths) {
  ::btool::app::collector::CollectStringsProperties(
      n, include_paths,
      [](const ::btool::node::PropertyStore *ps)
          -> const std::vector<std::string> * {
        return Properties::IncludePaths(ps);
      });
}

};  // namespace btool::app::collector::cc
