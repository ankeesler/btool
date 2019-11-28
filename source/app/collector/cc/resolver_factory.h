#ifndef BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORY_H_
#define BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORY_H_

#include <string>
#include <vector>

#include "node/node.h"

namespace btool::app::collector::cc {

class ResolverFactory {
 public:
  virtual class ::btool::node::Node::Resolver *NewCompileC(
      const std::vector<std::string> &include_dirs,
      const std::vector<std::string> &flags) = 0;
  virtual class ::btool::node::Node::Resolver *NewCompileCC(
      const std::vector<std::string> &include_dirs,
      const std::vector<std::string> &flags) = 0;
  virtual class ::btool::node::Node::Resolver *NewArchive() = 0;
  virtual class ::btool::node::Node::Resolver *NewLinkC(
      const std::vector<std::string> &flags) = 0;
  virtual class ::btool::node::Node::Resolver *NewLinkCC(
      const std::vector<std::string> &flags) = 0;
};

};  // namespace btool::app::collector::cc

#endif  // BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORY_H_
