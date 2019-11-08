#include "app/builder/builder.h"

#include <string>

#include "core/err.h"
#include "core/log.h"
#include "node/node.h"

namespace btool::app::builder {

::btool::core::VoidErr Builder::Build(const ::btool::node::Node &node) {
  ::btool::core::VoidErr err;

  node.Visit([&](const ::btool::node::Node *n) {
    if (!err) {
      auto current_err = c_->Current(*n);
      DEBUG("builder visiting %s, current: %s, resolver = %s\n",
            n->name().c_str(),
            (current_err ? current_err.Msg()
                         : (current_err.Ret() ? "true" : "false")),
            (n->resolver() == nullptr ? "null" : "something"));
      if (current_err) {
        err = ::btool::core::VoidErr::Failure(current_err.Msg());
      } else if (!current_err.Ret() && n->resolver() != nullptr) {
        err = n->resolver()->Resolve(*n);
      }
    }
  });

  return err;
}

};  // namespace btool::app::builder
