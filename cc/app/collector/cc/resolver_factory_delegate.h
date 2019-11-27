#ifndef BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORYDELEGATE_H_
#define BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORYDELEGATE_H_

#include <string>

#include "app/collector/cc/resolver_factory.h"
#include "app/collector/resolver_factory_delegate.h"
#include "node/node.h"
#include "node/property_store.h"

namespace btool::app::collector::cc {

class ResolverFactoryDelegate
    : public ::btool::app::collector::ResolverFactoryDelegate {
 public:
  ResolverFactoryDelegate(ResolverFactory *rf) : rf_(rf) {}
  ::btool::node::Node::Resolver *NewResolver(
      const std::string &name, const ::btool::node::PropertyStore &config,
      const ::btool::node::Node &n) override;

 private:
  ResolverFactory *rf_;
};

};  // namespace btool::app::collector::cc

#endif  // BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORYDELEGATE_H_
