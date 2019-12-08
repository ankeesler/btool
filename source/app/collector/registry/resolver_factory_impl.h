#ifndef BTOOL_APP_COLLECTOR_REGISTRY_RESOLVERFACTORYIMPL_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_RESOLVERFACTORYIMPL_H_

#include <string>
#include <vector>

#include "app/collector/registry/resolver_factory_delegate.h"
#include "node/node.h"

namespace btool::app::collector::registry {

class ResolverFactoryImpl : public ResolverFactoryDelegate::ResolverFactory {
 public:
  ~ResolverFactoryImpl() {
    for (auto a : allocations_) {
      delete a;
    }
  }

  ::btool::node::Node::Resolver *NewDownload(
      const std::string &url, const std::string &sha256) override;
  ::btool::node::Node::Resolver *NewUnzip() override;
  ::btool::node::Node::Resolver *NewUntar() override;
  ::btool::node::Node::Resolver *NewCmd(const std::string &path,
                                        const std::vector<std::string> &args,
                                        const std::string &dir) override;

 private:
  std::vector<::btool::node::Node::Resolver *> allocations_;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_RESOLVERFACTORYIMPL_H_
