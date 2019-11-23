#ifndef BTOOL_APP_COLLECTOR_REGISTRY_REGISTRYCOLLECTINI_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_REGISTRYCOLLECTINI_H_

#include "app/collector/base_collectini.h"
#include "app/collector/registry/registry.h"
#include "app/collector/store.h"

namespace btool::app::collector::registry {

class RegistryCollectini : public ::btool::app::collector::BaseCollectini {
 public:
  class GaggleCollector {
   public:
    virtual void Collect(::btool::app::collector::Store *s, Gaggle *g,
                         std::string root) = 0;
  };

  RegistryCollectini(Registry *r, std::string cache, GaggleCollector *gc)
      : r_(r), cache_(cache), gc_(gc) {}

  void Collect(::btool::app::collector::Store *s) override;

 private:
  Registry *r_;
  std::string cache_;
  GaggleCollector *gc_;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_REGISTRYCOLLECTINI_H_
