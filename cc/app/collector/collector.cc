#include "collector.h"

#include "app/collector/store.h"
#include "core/err.h"
#include "node/node.h"

namespace btool::app::collector {

::btool::core::VoidErr Collector::Collect() {
  Store s;
  for (auto c : cs_) {
    auto err = c->Collect(&s);
    if (err) {
      return err;
    }
  }
  return ::btool::core::VoidErr::Success();
}

};  // namespace btool::app::collector
