#include "app/builder/builder.h"

#include <string>

#include "core/err.h"
#include "node/node.h"

namespace btool::app::builder {

::btool::core::VoidErr Builder::Build(const ::btool::node::Node &node) {
  ::btool::core::VoidErr err;

  node.Visit([&](const ::btool::node::Node *n) {
    if (!err) {
      auto current_err = c_->Current(*n);
      if (current_err) {
        err = ::btool::core::VoidErr::Failure(current_err.Msg());
      } else if (!current_err.Ret() && n->Resolver() != nullptr) {
        err = n->Resolver()->Resolve(*n);
      }
    }
  });

  return err;
}

};  // namespace btool::app::builder
