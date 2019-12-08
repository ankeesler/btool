#include "app/collector/registry/gaggle_collector_impl.h"

#include "app/collector/registry/registry.h"
#include "app/collector/store.h"
#include "err.h"
#include "log.h"
#include "util/fs/fs.h"

namespace btool::app::collector::registry {

static void PrependRoot(::btool::node::Node *n, const std::string &root);

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

    PrependRoot(n_n, root);  // this is bad! see TODO below!

    for (auto rfd : rfds_) {
      auto r = rfd->NewResolver(n.resolver.name, n.resolver.config, *n_n);
      if (r != nullptr) {
        n_n->set_resolver(r);
        break;
      }
    }
    if (n_n->resolver() == nullptr) {
      DEBUGS() << "no known resolver for node " << n.name << " with resolver "
               << n.resolver.name;
    }
  }
}

static void PrependRoot(::btool::node::Node *n, const std::string &root) {
  // TODO: this is bad! we don't want to reach across to the cc package!
  auto prepend_root = [&root](std::string *s) { s->insert(0, root + "/"); };
  n->property_store()->ForEach("io.btool.collector.cc.includePaths",
                               prepend_root);
  n->property_store()->ForEach("io.btool.collector.cc.libraries", prepend_root);
}

};  // namespace btool::app::collector::registry
