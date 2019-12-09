#include "app/collector/registry/resolver_factory_delegate.h"

#include <string>
#include <vector>

#include "app/collector/properties.h"
#include "err.h"
#include "node/node.h"
#include "node/property_store.h"
#include "util/fs/fs.h"

namespace btool::app::collector::registry {

::btool::node::Node::Resolver *ResolverFactoryDelegate::NewResolver(
    const std::string &name, const ::btool::node::PropertyStore &config,
    const ::btool::node::Node &n) {
  if (name == "io.btool.collector.registry.resolvers/unzip") {
    return rf_->NewUnzip();
  } else if (name == "io.btool.collector.registry.resolvers/untar") {
    return rf_->NewUntar();
  } else if (name == "io.btool.collector.registry.resolvers/download") {
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
  } else if (name == "io.btool.collector.registry.resolvers/cmd") {
    const std::string *path = nullptr;
    const std::vector<std::string> *args = nullptr;
    const std::string *dir = nullptr;

    config.Read("path", &path);
    if (path == nullptr) {
      THROW_ERR("missing path key in cmd resolver config for node " + n.name());
    }

    config.Read("args", &args);

    std::string root =
        ::btool::app::collector::Properties::Root(n.property_store());
    std::string actual_dir;
    config.Read("dir", &dir);
    if (dir == nullptr) {
      actual_dir = root;
    } else {
      actual_dir = ::btool::util::fs::Join(root, *dir);
    }

    return rf_->NewCmd(*path,
                       args == nullptr ? std::vector<std::string>() : *args,
                       actual_dir);
  } else {
    return nullptr;
  }
}

};  // namespace btool::app::collector::registry
