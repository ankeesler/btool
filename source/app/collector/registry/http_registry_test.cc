#include "app/collector/registry/http_registry.h"

#include "gtest/gtest.h"

#include "app/collector/registry/registry.h"
#include "app/collector/registry/yaml_serializer.h"

TEST(HttpRegistry, DoesItFail) {
  ::btool::app::collector::registry::YamlSerializer ys;
  ::btool::app::collector::registry::HttpRegistry hr(
      "https://btoolregistry.cfapps.io", &ys);

  ::btool::app::collector::registry::Index i;
  hr.GetIndex(&i);

  ::btool::app::collector::registry::Gaggle g;
  hr.GetGaggle(i.files[0].path, &g);
}
