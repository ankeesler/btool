#ifndef BTOOL_APP_COLLECTOR_RESOLVERFACTORYIMPL_H_
#define BTOOL_APP_COLLECTOR_RESOLVERFACTORYIMPL_H_

#include <string>

#include "app/collector/resolver_factory.h"
#include "node/node.h"

namespace btool::app::collector {

class ResolverFactoryImpl : public ResolverFactory {
 public:
  ::btool::node::Node::Resolver *NewDownload(
      const std::string &url, const std::string &sha256) override;
  ::btool::node::Node::Resolver *NewUnzip() override;
};

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_RESOLVERFACTORY_H_
