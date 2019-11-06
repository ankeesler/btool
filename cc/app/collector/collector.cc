#include "collector.h"

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

::btool::core::VoidErr Collector::Collect() {
  Store s;
  StoreLogger sl;
  s.Listen(&sl);

  for (auto c : cs_) {
    auto err = c->Collect(&s);
    if (err) {
      return err;
    }
  }

  return ::btool::core::VoidErr::Success();
}

};  // namespace btool::app::collector
