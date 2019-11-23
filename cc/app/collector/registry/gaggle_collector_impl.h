#ifndef BTOOL_APP_COLLECTOR_REGISTRY_GAGGLECOLLECTORIMPL_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_GAGGLECOLLECTORIMPL_H_

#include <string>

#include "app/collector/registry/registry_collectini.h"
#include "app/collector/store.h"

namespace btool::app::collector::registry {

class GaggleCollectorImpl : RegistryCollectini::GaggleCollector {
 public:
  void Collect(::btool::app::collector::Store *s, Gaggle *g,
               std::string root) override;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_GAGGLECOLLECTORIMPL_H_
