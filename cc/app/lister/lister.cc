#include "app/lister/lister.h"

#include "node/node.h"

namespace btool::app::lister {

void Lister::List(const ::btool::node::Node &node) { node.String(os_); }

};  // namespace btool::app::lister
