#include "app/collector/registry/gaggle_collector_impl.h"

#include "app/collector/registry/registry.h"
#include "app/collector/store.h"
#include "err.h"
#include "util/fs/fs.h"

namespace btool::app::collector::registry {

void GaggleCollectorImpl::Collect(::btool::app::collector::Store *s, Gaggle *g,
                                  std::string root) {
  for (const auto &n : g->nodes) {
    auto name = ::btool::util::fs::Join(root, n.name);
    auto n_n = s->Put(name);

    for (const auto &d : n.dependencies) {
      if (d == "$this") {
        // TODO: handle me!
        continue;
      }

      name = ::btool::util::fs::Join(root, d);
      auto d_n = s->Get(name);
      if (d_n == nullptr) {
        THROW_ERR("unknown dependency " + d + " for node " + n.name);
      }
      n_n->dependencies()->push_back(d_n);
    }

    n_n->set_property_store(n.labels);

    for (auto rfd : rfds_) {
      auto r = rfd->NewResolver(n.resolver.name, n.resolver.config, root, *n_n);
      if (r != nullptr) {
        n_n->set_resolver(r);
        break;
      }
    }
    if (n_n->resolver() == nullptr) {
      THROW_ERR("no known resolver for node " + n.name + " with resolver " +
                n.resolver.name);
    }
  }
}

};  // namespace btool::app::collector::registry
