#ifndef BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORY_H_
#define BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORY_H_

#include <string>
#include <vector>

#include "node/node.h"

namespace btool::app::collector::cc {

class ResolverFactory {
 public:
  virtual class ::btool::node::Node::Resolver *NewCompileC() = 0;
  virtual class ::btool::node::Node::Resolver *NewCompileCC() = 0;
  virtual class ::btool::node::Node::Resolver *NewArchive() = 0;
  virtual class ::btool::node::Node::Resolver *NewLinkC() = 0;
  virtual class ::btool::node::Node::Resolver *NewLinkCC() = 0;
};

};  // namespace btool::app::collector::cc

#endif  // BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORY_H_
