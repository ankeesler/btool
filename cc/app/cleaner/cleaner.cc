#include "app/cleaner/cleaner.h"

#include <string>

#include "core/err.h"
#include "node/node.h"

namespace btool::app::cleaner {

::btool::core::VoidErr Cleaner::Clean(const ::btool::node::Node& node) {
  ::btool::core::VoidErr err;

  node.Visit([&](const ::btool::node::Node* n) {
    if (!err) {
      err = ra_->RemoveAll(n->name());
    }
  });

  return err;
}

};  // namespace btool::app::cleaner
