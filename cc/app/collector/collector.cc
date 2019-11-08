#include "collector.h"

#include <string>

#include "app/collector/store.h"
#include "core/err.h"
#include "core/log.h"
#include "node/node.h"

namespace btool::app::collector {

::btool::core::Err<::btool::node::Node *> Collector::Collect(
    const std::string &target) {
  for (auto c : cs_) {
    c->Collect(s_);
  }

  auto n = s_->Get(target);
  if (n == nullptr) {
    return ::btool::core::Err<::btool::node::Node *>::Failure("unknown target");
  }

  return ::btool::core::Err<::btool::node::Node *>::Success(n);
}

};  // namespace btool::app::collector
