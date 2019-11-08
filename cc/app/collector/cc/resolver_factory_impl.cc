#include "app/collector/cc/resolver_factory_impl.h"

#include <string>
#include <vector>

#include "node/node.h"

namespace btool::app::collector::cc {

::btool::node::Node::Resolver *ResolverFactoryImpl::NewCompileC(
    const std::vector<std::string> &include_dirs,
    const std::vector<std::string> &flags) {
  return nullptr;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewCompileCC(
    const std::vector<std::string> &include_dirs,
    const std::vector<std::string> &flags) {
  return nullptr;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewArchive() {
  return nullptr;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewLinkC(
    const std::vector<std::string> &flags) {
  return nullptr;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewLinkCC(
    const std::vector<std::string> &flags) {
  return nullptr;
}

};  // namespace btool::app::collector::cc
