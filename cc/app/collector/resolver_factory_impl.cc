#include "app/collector/resolver_factory_impl.h"

#include <string>

#include "node/node.h"

namespace btool::app::collector {

::btool::node::Node::Resolver *ResolverFactoryImpl::NewDownload(
    const std::string &url, const std::string &sha256) {
  return nullptr;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewUnzip() {
  return nullptr;
}

};  // namespace btool::app::collector
