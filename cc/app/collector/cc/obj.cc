#include "obj.h"

#include <string>
#include <vector>

#include "app/collector/cc/properties.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "node/node.h"
#include "util/string/string.h"

namespace btool::app::collector::cc {

void CollectIncludePaths(::btool::node::Node *n,
                         std::vector<std::string> *include_paths);

void Obj::OnSet(::btool::app::collector::Store *s, const std::string &name) {
  auto d = s->Get(name);
  if (d == nullptr) {
    return;
  }

  if (!::btool::app::collector::Properties::Local(d->property_store())) {
    return;
  }

  bool c = ::btool::util::string::HasSuffix(name.c_str(), ".c");
  bool cc = ::btool::util::string::HasSuffix(name.c_str(), ".cc");
  if (!c && !cc) {
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
  n->SetResolver(r);
  n->dependencies()->push_back(d);
}

void CollectIncludePaths(::btool::node::Node *n,
                         std::vector<std::string> *include_paths) {
  n->Visit([include_paths](const ::btool::node::Node *vn) {
    auto ips = Properties::IncludePaths(vn->property_store());
    if (ips != nullptr) {
      include_paths->insert(include_paths->end(), ips->begin(), ips->end());
    }
  });
}

};  // namespace btool::app::collector::cc
