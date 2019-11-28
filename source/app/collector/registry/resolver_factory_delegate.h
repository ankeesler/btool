#ifndef BTOOL_APP_COLLECTOR_REGISTRY_RESOLVERFACTORYDELEGATE_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_RESOLVERFACTORYDELEGATE_H_

#include <string>

#include "app/collector/resolver_factory_delegate.h"
#include "node/node.h"
#include "node/property_store.h"

namespace btool::app::collector::registry {

class ResolverFactoryDelegate
    : public ::btool::app::collector::ResolverFactoryDelegate {
 public:
  class ResolverFactory {
   public:
    virtual ::btool::node::Node::Resolver *NewDownload(
        const std::string &url, const std::string &sha256) = 0;
    virtual ::btool::node::Node::Resolver *NewUnzip() = 0;
    virtual ::btool::node::Node::Resolver *NewUntar() = 0;
  };

  ResolverFactoryDelegate(ResolverFactory *rf) : rf_(rf) {}
  ::btool::node::Node::Resolver *NewResolver(
      const std::string &name, const ::btool::node::PropertyStore &config,
      const ::btool::node::Node &n) override;

 private:
  ResolverFactory *rf_;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_RESOLVERFACTORYDELEGATE_H_
