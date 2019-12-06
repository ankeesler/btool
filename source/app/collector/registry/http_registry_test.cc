#include "app/collector/registry/http_registry.h"

#include "gtest/gtest.h"

#include "app/collector/registry/registry.h"
#include "app/collector/registry/yaml_serializer.h"
#include "util/fs/fs.h"

TEST(HttpRegistry, DoesItFail) {
  auto dir = ::btool::util::fs::TempDir();

  ::btool::app::collector::registry::YamlSerializer<
      ::btool::app::collector::registry::Index>
      ys_i;
  ::btool::app::collector::registry::YamlSerializer<
      ::btool::app::collector::registry::Gaggle>
      ys_g;
  ::btool::app::collector::registry::HttpRegistry hr(
      "https://btoolregistry.cfapps.io", &ys_i, &ys_g);

  ::btool::app::collector::registry::Index i;
  hr.GetIndex(&i);

  ::btool::app::collector::registry::Gaggle g;
  hr.GetGaggle(i.files[0].path, &g);

  ::btool::util::fs::RemoveAll(dir);
}
