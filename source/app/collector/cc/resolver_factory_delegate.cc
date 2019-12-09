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
    std::vector<std::string> include_paths;
    ::btool::app::collector::CollectStringsProperties(
        n, &include_paths, [](const ::btool::node::PropertyStore *ps) {
          return Properties::IncludePaths(ps);
        });
    return rf_->NewCompileC(include_paths, {});
  } else if (name == "io.btool.collector.cc.resolvers/compileCC") {
    std::vector<std::string> include_paths;
    ::btool::app::collector::CollectStringsProperties(
        n, &include_paths, [](const ::btool::node::PropertyStore *ps) {
          return Properties::IncludePaths(ps);
        });
    return rf_->NewCompileCC(include_paths, {});
  } else if (name == "io.btool.collector.cc.resolvers/archive") {
    return rf_->NewArchive();
  } else if (name == "io.btool.collector.cc.resolvers/linkC") {
    std::vector<std::string> flags;
    ::btool::app::collector::CollectStringsProperties(
        n, &flags, [](const ::btool::node::PropertyStore *ps) {
          return Properties::LinkFlags(ps);
        });
    return rf_->NewLinkC(flags);
  } else if (name == "io.btool.collector.cc.resolvers/linkCC") {
    std::vector<std::string> flags;
    ::btool::app::collector::CollectStringsProperties(
        n, &flags, [](const ::btool::node::PropertyStore *ps) {
          return Properties::LinkFlags(ps);
        });
    return rf_->NewLinkCC(flags);
  } else {
    return nullptr;
  }
}  // namespace btool::app::collector::cc
};  // namespace btool::app::collector::cc
