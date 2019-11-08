#ifndef BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORYIMPL_H_
#define BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORYIMPL_H_

#include <string>
#include <vector>

#include "app/collector/cc/resolver_factory.h"
#include "node/node.h"

namespace btool::app::collector::cc {

class ResolverFactoryImpl : public ResolverFactory {
 public:
  ::btool::node::Node::Resolver *NewCompileC(
      const std::vector<std::string> &include_dirs,
      const std::vector<std::string> &flags) override;
  ::btool::node::Node::Resolver *NewCompileCC(
      const std::vector<std::string> &include_dirs,
      const std::vector<std::string> &flags) override;
  ::btool::node::Node::Resolver *NewArchive() override;
  ::btool::node::Node::Resolver *NewLinkC(
      const std::vector<std::string> &flags) override;
  ::btool::node::Node::Resolver *NewLinkCC(
      const std::vector<std::string> &flags) override;
};

};  // namespace btool::app::collector::cc

#endif  // BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORYIMPL_H_
