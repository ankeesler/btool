#include "app/cleaner/cleaner.h"

#include <string>

#include "node/node.h"

namespace btool::app::cleaner {

bool Cleaner::Clean(const ::btool::node::Node& node, std::string* err) {
  bool success = true;

  node.Visit([&](const ::btool::node::Node* n) {
    if (success) {
      success = ra_->RemoveAll(n->Name(), err);
    }
  });

  return success;
}

};  // namespace btool::app::cleaner
