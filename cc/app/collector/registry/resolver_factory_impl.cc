#include "app/collector/registry/resolver_factory_impl.h"

#include "app/collector/registry/registry.h"
#include "node/node.h"

namespace btool::app::collector::registry {

::btool::node::Node::Resolver *ResolverFactoryImpl::New(const Resolver &r) {
  return nullptr;
}

};  // namespace btool::app::collector::registry
