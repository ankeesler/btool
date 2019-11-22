#include "app/lister/lister.h"

#include "node/node.h"

namespace btool::app::lister {

bool Lister::List(const ::btool::node::Node &node, std::string *ret_err) {
  node.String(os_);
  return true;
}

};  // namespace btool::app::lister
