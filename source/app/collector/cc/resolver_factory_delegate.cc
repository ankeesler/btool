#include "app/collector/cc/resolver_factory_delegate.h"

#include <string>
#include <vector>

#include "app/collector/cc/properties.h"
#include "app/collector/properties.h"
#include "err.h"
#include "node/node.h"
#include "node/property_store.h"

namespace btool::app::collector::cc {

::btool::node::Node::Resolver *ResolverFactoryDelegate::NewResolver(
    const std::string &name, const ::btool::node::PropertyStore &config,
    const ::btool::node::Node &n) {
  if (name == "io.btool.collector.cc.resolvers/compileC") {
    return rf_->NewCompileC();
  } else if (name == "io.btool.collector.cc.resolvers/compileCC") {
    return rf_->NewCompileCC();
  } else if (name == "io.btool.collector.cc.resolvers/archive") {
    return rf_->NewArchive();
  } else if (name == "io.btool.collector.cc.resolvers/linkC") {
    return rf_->NewLinkC();
  } else if (name == "io.btool.collector.cc.resolvers/linkCC") {
    return rf_->NewLinkCC();
  } else {
    return nullptr;
  }
}  // namespace btool::app::collector::cc
};  // namespace btool::app::collector::cc
