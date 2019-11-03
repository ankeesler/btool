#include "app/collector/resolver_factory_impl.h"

#include "node/node.h"

namespace btool::app::collector {

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

::btool::node::Node::Resolver *ResolverFactoryImpl::NewDownload(
    const std::string &url, const std::string &sha256) {
  return nullptr;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewUnzip() {
  return nullptr;
}

};  // namespace btool::app::collector
