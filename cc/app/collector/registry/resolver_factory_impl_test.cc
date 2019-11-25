#include "app/collector/registry/resolver_factory_impl.h"

#include <memory>

#include "gtest/gtest.h"

#include "app/collector/registry/registry.h"
#include "node/node.h"

MockResolverFactoryDelegatere

TEST(ResolverFactoryImpl, CompileC) {
  ::btool::app::collector::registry::ResolverFactoryImpl rfi;

  ::btool::node::Node dd("dd");
  dd.property_store()->Append("io.btool.collector.cc.includePaths",
                              "dd/include");
  dd.property_store()->Append("io.btool.collector.cc.includePaths",
                              "dd/source");
  dd.property_store()->Append("io.btool.collector.cc.linkFlags", "-ldd");

  ::btool::node::Node d("d");
  d->dependencies()->push_back(&dd);
  d.property_store()->Append("io.btool.collector.cc.includePaths", "d/include");
  d.property_store()->Append("io.btool.collector.cc.includePaths", "d/source");
  dd.property_store()->Append("io.btool.collector.cc.linkFlags", "-ld");

  ::btool::node::Node n("n");
  n->dependencies()->push_back(&d);
  n.property_store()->Append("io.btool.collector.cc.includePaths", "source");
  n.property_store()->Append("io.btool.collector.cc.linkFlags", "-lpthread");

  rfi.New(&n);
}
