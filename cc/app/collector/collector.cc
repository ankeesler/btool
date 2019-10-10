#include "app/collector/collector.h"

#include "core/err.h"
#include "node/node.h"
#include "node/store.h"

namespace btool::app::collector {

::btool::core::VoidErr Collector::Collect() {
  ::btool::node::Store s;
  for (auto c : cs_) {
    auto err = c->Collect(&s);
    if (err) {
      return err;
    }
  }
  return ::btool::core::VoidErr::Success();
}

};  // namespace btool::app::collector
