#include "collector.h"

#include <string>

#include "app/collector/store.h"
#include "err.h"
#include "log.h"
#include "node/node.h"

namespace btool::app::collector {

::btool::Err<::btool::node::Node *> Collector::Collect(
    const std::string &target) {
  for (auto c : cs_) {
    c->Collect(s_);
  }

  auto n = s_->Get(target);
  if (n == nullptr) {
    return ::btool::Err<::btool::node::Node *>::Failure("unknown target");
  }

  return ::btool::Err<::btool::node::Node *>::Success(n);
}

};  // namespace btool::app::collector
