#ifndef BTOOL_APP_COLLECTOR_RESOLVERFACTORY_H_
#define BTOOL_APP_COLLECTOR_RESOLVERFACTORY_H_

#include <string>

#include "node/node.h"

namespace btool::app::collector {

class ResolverFactory {
 public:
  virtual class ::btool::node::Node::Resolver *NewDownload(
      const std::string &url, const std::string &sha256) = 0;
  virtual class ::btool::node::Node::Resolver *NewUnzip() = 0;
};

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_RESOLVERFACTORY_H_
