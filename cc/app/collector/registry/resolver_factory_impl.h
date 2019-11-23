#ifndef BTOOL_APP_COLLECTOR_REGISTRY_RESOLVERFACTORYIMPL_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_RESOLVERFACTORYIMPL_H_

#include "app/collector/registry/gaggle_collector_impl.h"
#include "app/collector/registry/registry.h"
#include "node/node.h"

namespace btool::app::collector::registry {

class ResolverFactoryImpl : public GaggleCollectorImpl::ResolverFactory {
 public:
  ::btool::node::Node::Resolver *New(const Resolver &r) override;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_RESOLVERFACTORYIMPL_H_
