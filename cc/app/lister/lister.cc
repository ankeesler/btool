#include "app/lister/lister.h"

#include "core/err.h"
#include "node/node.h"

namespace btool::app::lister {

::btool::core::VoidErr Lister::List(const ::btool::node::Node &node) {
  node.String(os_);
  return ::btool::core::VoidErr::Success();
}

};  // namespace btool::app::lister
