#ifndef BTOOL_APP_COLLECTOR_REGISTRY_GAGGLECOLLECTORIMPL_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_GAGGLECOLLECTORIMPL_H_

#include <string>

#include "app/collector/registry/registry_collectini.h"
#include "app/collector/store.h"
#include "node/node.h"

namespace btool::app::collector::registry {

class GaggleCollectorImpl : public RegistryCollectini::GaggleCollector {
 public:
  class ResolverFactory {
   public:
    virtual ::btool::node::Node::Resolver *New(const Resolver &r) = 0;
  };

  GaggleCollectorImpl(ResolverFactory *rf) : rf_(rf) {}

  void Collect(::btool::app::collector::Store *s, Gaggle *g,
               std::string root) override;

 private:
  ResolverFactory *rf_;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_GAGGLECOLLECTORIMPL_H_
