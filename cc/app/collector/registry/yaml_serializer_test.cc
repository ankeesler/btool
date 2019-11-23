#include "app/collector/registry/yaml_serializer.h"

#include <sstream>

#include "gtest/gtest.h"

TEST(YamlSerializer, UnmarshalIndex) {
  ::btool::app::collector::registry::IndexFile if0{.path = "some/path0.yml",
                                                   .sha256 = "sha0"};
  ::btool::app::collector::registry::IndexFile if1{.path = "some/path1.yml",
                                                   .sha256 = "sha1"};
  ::btool::app::collector::registry::Index ex_i{.files = {if0, if1}};

  std::string content =
      "---\n"
      "- path: some/path0.yml\n"
      "  sha256: sha0\n"
      "- path: some/path1.yml\n"
      "  sha256: sha1\n";
  std::stringstream ss{content};
  ::btool::app::collector::registry::YamlSerializer ys;

  ::btool::app::collector::registry::Index ac_i;
  ys.UnmarshalIndex(&ss, &ac_i);
  EXPECT_EQ(ex_i, ac_i);
}
