#ifndef BTOOL_APP_COLLECTOR_REGISTRY_GAGGLECOLLECTORIMPL_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_GAGGLECOLLECTORIMPL_H_

#include <string>
#include <vector>

#include "app/collector/registry/registry_collectini.h"
#include "app/collector/resolver_factory_delegate.h"
#include "app/collector/store.h"
#include "node/node.h"

namespace btool::app::collector::registry {

class GaggleCollectorImpl : public RegistryCollectini::GaggleCollector {
 public:
  class ResolverFactory {
   public:
    virtual ::btool::node::Node::Resolver *New(const Resolver &r) = 0;
  };

  void AddResolverFactoryDelegate(
      ::btool::app::collector::ResolverFactoryDelegate *rfd) {
    rfds_.push_back(rfd);
  }
  void Collect(::btool::app::collector::Store *s, Gaggle *g,
               std::string root) override;

 private:
  std::vector<::btool::app::collector::ResolverFactoryDelegate *> rfds_;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_GAGGLECOLLECTORIMPL_H_
