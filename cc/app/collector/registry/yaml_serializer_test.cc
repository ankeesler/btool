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
      "files:\n"
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

TEST(YamlSerializer, UnmarshalGaggle) {
  ::btool::node::PropertyStore rc;
  rc.Write("sup", "homie");
  ::btool::app::collector::registry::Resolver r = {.name = "hey", .config = rc};

  ::btool::node::PropertyStore ls;
  ls.Write("some-bool", true);
  ls.Append("some-strings", "this");
  ls.Append("some-strings", "is");
  ls.Append("some-strings", "a");
  ls.Append("some-strings", "list");
  ::btool::app::collector::registry::Node n0{
      .name = "tuna", .dependencies = {"fish", "marlin"}, .resolver = r};
  ::btool::app::collector::registry::Node n1{
      .name = "another-tuna",
      .dependencies = {"another-fish", "another-marlin"},
      .labels = ls,
      .resolver = r};

  ::btool::app::collector::registry::Gaggle ex_g{.nodes = {n0, n1}};

  std::string content =
      "---\n"
      "nodes:\n"
      "- name: tuna\n"
      "  dependencies: [fish, marlin]\n"
      "  resolver:\n"
      "    name: hey\n"
      "    config:\n"
      "      sup: homie\n"
      "- name: another-tuna\n"
      "  dependencies: [another-fish, another-marlin]\n"
      "  labels:\n"
      "    some-bool: true\n"
      "    some-strings: [this, is, a, list]\n"
      "  resolver:\n"
      "    name: hey\n"
      "    config:\n"
      "      sup: homie\n";
  std::stringstream ss{content};
  ::btool::app::collector::registry::YamlSerializer ys;

  ::btool::app::collector::registry::Gaggle ac_g;
  ys.UnmarshalGaggle(&ss, &ac_g);
  EXPECT_EQ(ex_g, ac_g);
}
