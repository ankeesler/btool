#ifndef BTOOL_APP_COLLECTOR_REGISTRY_REGISTRYCOLLECTINI_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_REGISTRYCOLLECTINI_H_

#include "app/collector/base_collectini.h"
#include "app/collector/registry/registry.h"
#include "app/collector/store.h"
#include "util/cache.h"

namespace btool::app::collector::registry {

class RegistryCollectini : public ::btool::app::collector::BaseCollectini {
 public:
  class GaggleCollector {
   public:
    virtual void Collect(::btool::app::collector::Store *s, Gaggle *g,
                         std::string root) = 0;
  };

  RegistryCollectini(Registry *r, std::string cache,
                     ::btool::util::Cache<Index> *c_i,
                     ::btool::util::Cache<Gaggle> *c_g, GaggleCollector *gc)
      : r_(r), cache_(cache), c_i_(c_i), c_g_(c_g), gc_(gc) {}

  void Collect(::btool::app::collector::Store *s) override;

 private:
  Registry *r_;
  std::string cache_;
  ::btool::util::Cache<Index> *c_i_;
  ::btool::util::Cache<Gaggle> *c_g_;
  GaggleCollector *gc_;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_REGISTRYCOLLECTINI_H_
