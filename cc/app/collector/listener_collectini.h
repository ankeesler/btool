#ifndef BTOOL_APP_COLLECTOR_LISTENERCOLLECTINI_H_
#define BTOOL_APP_COLLECTOR_LISTENERCOLLECTINI_H_

#include "app/collector/collector.h"
#include "app/collector/store.h"

namespace btool::app::collector {

class ListenerCollectini : public Store::Listener,
                           public Collector::Collectini {
 public:
  ::btool::core::VoidErr Collect(Store *s) override {
    s->Listen(this);
    return ::btool::core::VoidErr::Success();
  }
};

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_LISTENERCOLLECTINI_H_
