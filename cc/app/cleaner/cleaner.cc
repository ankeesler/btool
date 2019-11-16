#include "app/cleaner/cleaner.h"

#include <string>

#include "err.h"
#include "core/log.h"
#include "node/node.h"

namespace btool::app::cleaner {

::btool::VoidErr Cleaner::Clean(const ::btool::node::Node& node) {
  ::btool::VoidErr err;

  node.Visit([&](const ::btool::node::Node* n) {
    if (!err && n->resolver() != nullptr) {
      DEBUG("cleaning %s\n", n->name().c_str());
      err = ra_->RemoveAll(n->name());
    }
  });

  return err;
}

};  // namespace btool::app::cleaner
