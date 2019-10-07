#include "app/builder/builder.h"

#include <string>

#include "node/node.h"

namespace btool::app::builder {

bool Builder::Build(const ::btool::node::Node &node, std::string *err) {
  bool success = true;

  node.Visit([&](const ::btool::node::Node *n) {
    if (success) {
      bool current;
      success = c_->Current(*n, &current, err);
      if (!current) {
        success = n->Resolver()->Resolve(*n, err);
      }
    }
  });

  return success;
}

};  // namespace btool::app::builder
