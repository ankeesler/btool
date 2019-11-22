#include "app/cleaner/cleaner.h"

#include <string>

#include "err.h"
#include "log.h"
#include "node/node.h"

namespace btool::app::cleaner {

bool Cleaner::Clean(const ::btool::node::Node& node, std::string* ret_err) {
  bool success = true;

  node.Visit([&](const ::btool::node::Node* n) {
    if (success && n->resolver() != nullptr) {
      DEBUGS() << "cleaning " << n->name() << std::endl;

      std::string err;
      if (!ra_->RemoveAll(n->name(), &err)) {
        *ret_err = ::btool::WrapErr(err, "remove all");
        success = false;
      }
    }
  });

  return success;
}

};  // namespace btool::app::cleaner
