#include "app/collector/registry/resolver_factory_delegate.h"

#include <string>

#include "err.h"
#include "node/node.h"
#include "node/property_store.h"

namespace btool::app::collector::registry {

::btool::node::Node::Resolver *ResolverFactoryDelegate::NewResolver(
    const std::string &name, const ::btool::node::PropertyStore &config,
    const ::btool::node::Node &n) {
  if (name == "unzip") {
    return rf_->NewUnzip();
  } else if (name == "download") {
    const std::string *url = nullptr;
    const std::string *sha256 = nullptr;

    config.Read("url", &url);
    if (url == nullptr) {
      THROW_ERR("missing url key in download resolver config for node " +
                n.name());
    }

    config.Read("sha256", &sha256);
    if (sha256 == nullptr) {
      THROW_ERR("missing sha256 key in download resolver config for node " +
                n.name());
    }

    return rf_->NewDownload(*url, *sha256);
  } else {
    return nullptr;
  }
}

};  // namespace btool::app::collector::registry
