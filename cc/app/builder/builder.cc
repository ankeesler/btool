#include "app/builder/builder.h"

#include <string>

#include "err.h"
#include "log.h"
#include "node/node.h"

namespace btool::app::builder {

void Builder::Build(const ::btool::node::Node &node) {
  node.Visit([&](const ::btool::node::Node *n) {
    bool current = c_->Current(*n);
    DEBUG("builder visiting %s, current: %s, resolver = %s\n",
          n->name().c_str(), (current ? "true" : "false"),
          (n->resolver() == nullptr ? "null" : "something"));
    if (!current && n->resolver() != nullptr) {
      n->resolver()->Resolve(*n);
    }
  });
}

}  // namespace btool::app::builder
