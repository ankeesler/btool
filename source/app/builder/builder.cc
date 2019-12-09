#include "app/builder/builder.h"

#include <string>

#include "err.h"
#include "log.h"
#include "node/node.h"

namespace btool::app::builder {

void Builder::Build(const ::btool::node::Node &node) {
  node.Visit([&](const ::btool::node::Node *n) {
    bool current = cu_->Current(*n);
    DEBUG("builder visiting %s, current: %s, resolver = %s\n",
          n->name().c_str(), (current ? "true" : "false"),
          (n->resolver() == nullptr ? "null" : "something"));
    if (n->resolver() != nullptr) {
      ca_->OnPreResolve(*n, current);
      if (!current) {
        n->resolver()->Resolve(*n);
      }
      ca_->OnPostResolve(*n, current);
    }
  });
}

}  // namespace btool::app::builder
