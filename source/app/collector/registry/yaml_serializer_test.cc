#include "app/collector/registry/yaml_serializer.h"

#include <iostream>
#include <sstream>

#include "gtest/gtest.h"

#include "app/collector/registry/registry.h"
#include "log.h"

TEST(YamlSerializer, Index) {
  ::btool::app::collector::registry::IndexFile if0{.path = "some/path0.yml",
                                                   .sha256 = "sha0"};
  ::btool::app::collector::registry::IndexFile if1{.path = "some/path1.yml",
                                                   .sha256 = "sha1"};
  ::btool::app::collector::registry::Index ex_i{.files = {if0, if1}};

  std::string content =
      "files:\n"
      "  - path: some/path0.yml\n"
      "    sha256: sha0\n"
      "  - path: some/path1.yml\n"
      "    sha256: sha1";
  std::stringstream ex_ss{content};
  ::btool::app::collector::registry::YamlSerializer<
      ::btool::app::collector::registry::Index>
      ys;

  ::btool::app::collector::registry::Index ac_i;
  ys.Unmarshal(&ex_ss, &ac_i);
  EXPECT_EQ(ex_i, ac_i);

  std::stringstream ac_ss;
  ys.Marshal(&ac_ss, ex_i);
  EXPECT_EQ(ex_ss.str(), ac_ss.str());
}

TEST(YamlSerializer, Gaggle) {
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
      "nodes:\n"
      "  - name: tuna\n"
      "    dependencies:\n"
      "      - fish\n"
      "      - marlin\n"
      "    labels: ~\n"
      "    resolver:\n"
      "      name: hey\n"
      "      config:\n"
      "        sup: homie\n"
      "  - name: another-tuna\n"
      "    dependencies:\n"
      "      - another-fish\n"
      "      - another-marlin\n"
      "    labels:\n"
      "      some-bool: true\n"
      "      some-strings:\n"
      "        - this\n"
      "        - is\n"
      "        - a\n"
      "        - list\n"
      "    resolver:\n"
      "      name: hey\n"
      "      config:\n"
      "        sup: homie";
  std::stringstream ex_ss{content};
  ::btool::app::collector::registry::YamlSerializer<
      ::btool::app::collector::registry::Gaggle>
      ys;

  ::btool::app::collector::registry::Gaggle ac_g;
  ys.Unmarshal(&ex_ss, &ac_g);
  EXPECT_EQ(ex_g, ac_g);

  std::stringstream ac_ss;
  ys.Marshal(&ac_ss, ex_g);
  EXPECT_EQ(ex_ss.str(), ac_ss.str());
}
