#include "app/lister/lister.h"

#include "err.h"
#include "node/node.h"

namespace btool::app::lister {

::btool::VoidErr Lister::List(const ::btool::node::Node &node) {
  node.String(os_);
  return ::btool::VoidErr::Success();
}

};  // namespace btool::app::lister
