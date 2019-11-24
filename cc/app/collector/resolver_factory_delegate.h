#ifndef BTOOL_APP_COLLECTOR_RESOLVERFACTORYDELEGATE_H_
#define BTOOL_APP_COLLECTOR_RESOLVERFACTORYDELEGATE_H_

#include <string>

#include "node/node.h"
#include "node/property_store.h"

namespace btool::app::collector {

class ResolverFactoryDelegate {
 public:
  virtual ::btool::node::Node::Resolver *NewResolver(
      const std::string &name, const ::btool::node::PropertyStore &config,
      const std::string &root, const ::btool::node::Node &n) = 0;
};

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_RESOLVERFACTORYDELEGATE_H_
