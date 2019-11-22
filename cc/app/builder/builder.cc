#include "app/builder/builder.h"

#include <string>

#include "err.h"
#include "log.h"
#include "node/node.h"

namespace btool::app::builder {

bool Builder::Build(const ::btool::node::Node &node, std::string *ret_err) {
  bool success = true;

  node.Visit([this, ret_err, &success](const ::btool::node::Node *n) {
    if (success) {
      bool ret_current;
      std::string current_ret_err;
      auto success = c_->Current(*n, &ret_current, &current_ret_err);
      DEBUGS() << "builer visiting " << n->name() << ", current: "
               << (success ? (ret_current ? "true" : "false") : current_ret_err)
               << ", resolver = "
               << (n->resolver() == nullptr ? "null" : "something")
               << std::endl;
      if (!success) {
        success = false;
        *ret_err = ::btool::WrapErr(current_ret_err, "current");
      } else if (!ret_current && n->resolver() != nullptr) {
        auto err = n->resolver()->Resolve(*n);
        if (err) {
          success = false;
          *ret_err = err.Msg();
        }
      }
    }
  });

  return success;
}

};  // namespace btool::app::builder
