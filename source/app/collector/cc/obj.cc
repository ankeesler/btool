#include "app/collector/cc/obj.h"

#include <iostream>
#include <string>
#include <vector>

#include "app/collector/cc/properties.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "log.h"
#include "node/node.h"
#include "util/string/string.h"

namespace btool::app::collector::cc {

#define DROP(name, reason) \
  DEBUGS() << "obj: drop " << (name) << " (" << reason << ")" << std::endl

void Obj::OnNotify(::btool::app::collector::Store *s, const std::string &name) {
  auto d = s->Get(name);
  if (d == nullptr) {
    DROP(name, "unknown name");
    return;
  }

  if (!::btool::app::collector::Properties::Local(d->property_store())) {
    DROP(name, "not local");
    return;
  }

  bool c = ::btool::util::string::HasSuffix(name, ".c");
  bool cc = ::btool::util::string::HasSuffix(name, ".cc");
  if (!c && !cc) {
    DROP(name, "not .c/.cc");
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

  DEBUGS() << "added new object " << obj_name << std::endl;

  ::btool::node::Node::Resolver *r;
  if (c) {
    r = rf_->NewCompileC();
  } else {  // cc
    r = rf_->NewCompileCC();
  }
  n->set_resolver(r);
  n->dependencies()->push_back(d);
  Notify(s, n->name());
}

};  // namespace btool::app::collector::cc
