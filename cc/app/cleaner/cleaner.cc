#include "app/cleaner/cleaner.h"

#include <string>

#include "err.h"
#include "log.h"
#include "node/node.h"

namespace btool::app::cleaner {

void Cleaner::Clean(const ::btool::node::Node& node) {
  node.Visit([&](const ::btool::node::Node* n) {
    if (n->resolver() != nullptr) {
      DEBUGS() << "cleaning " << n->name() << std::endl;
      ra_->RemoveAll(n->name());
    }
  });
}

};  // namespace btool::app::cleaner
