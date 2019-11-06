#include "collector.h"

#include <string>

#include "app/collector/store.h"
#include "core/err.h"
#include "core/log.h"
#include "node/node.h"

namespace btool::app::collector {

class StoreLogger : public Store::Listener {
 public:
  void OnSet(Store *s, const std::string &name) override {
    DEBUG("store logger: set %s\n", name.c_str());
  }
};

::btool::core::Err<::btool::node::Node *> Collector::Collect(
    const std::string &target) {
  StoreLogger sl;
  s_->Listen(&sl);

  for (auto c : cs_) {
    auto err = c->Collect(s_);
    if (err) {
      return ::btool::core::Err<::btool::node::Node *>::Failure(err.Msg());
    }
  }

  auto n = s_->Get(target);
  if (n == nullptr) {
    return ::btool::core::Err<::btool::node::Node *>::Failure("unknown target");
  }

  return ::btool::core::Err<::btool::node::Node *>::Success(n);
}

};  // namespace btool::app::collector
