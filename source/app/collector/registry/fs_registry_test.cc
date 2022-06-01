#include "app/collector/registry/fs_registry.h"

#include <string>

#include "app/collector/registry/testing/registry.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "node/node.h"
#include "util/fs/fs.h"

using ::testing::ContainerEq;
using ::testing::InSequence;
using ::testing::SetArgPointee;
using ::testing::StrictMock;

class FsRegistryTest : public ::testing::Test {
 protected:
  void SetUp() override {
    root_ = ::btool::util::fs::TempDir();

    auto index_yml = ::btool::util::fs::Join(root_, "index.yml");
    ::btool::util::fs::WriteFile(index_yml, "here is an index.yml file\n");

    auto g0 = ::btool::util::fs::Join(root_, "some-path-0");
    ::btool::util::fs::WriteFile(g0, "here is an gaggle 0 file\n");

    auto g1 = ::btool::util::fs::Join(root_, "some-path-1");
    ::btool::util::fs::WriteFile(g1, "here is an gaggle 1 file\n");
  }

  void TearDown() override { ::btool::util::fs::RemoveAll(root_); }

  std::string root_;
};

TEST_F(FsRegistryTest, Success) {
  InSequence s;

  ::btool::app::collector::registry::IndexFile if0{.path = "some-path-0",
                                                   .sha256 = "sha0"};
  ::btool::app::collector::registry::IndexFile if1{.path = "some-path-1",
                                                   .sha256 = "sha1"};
  ::btool::app::collector::registry::Index ex_i{.files = {if0, if1}};

  ::btool::node::PropertyStore ps;
  ps.Write("bool-property", true);
  ps.Append("strings-property", "some-string");
  ::btool::app::collector::registry::Resolver r = {.name = "r", .config = ps};
  ::btool::app::collector::registry::Node a{.name = "a", .labels = ps};
  ::btool::app::collector::registry::Node b{
      .name = "b", .dependencies = {"a"}, .resolver = r};
  ::btool::app::collector::registry::Node c{.name = "c", .dependencies = {"b"}};
  ::btool::app::collector::registry::Node d{.name = "d",
                                            .dependencies = {"b", "c"}};

  ::btool::app::collector::registry::Gaggle ex_g0{.nodes = {a, b}};
  ::btool::app::collector::registry::Gaggle ex_g1{.nodes = {c, d}};

  StrictMock<::btool::app::collector::registry::testing::MockSerializer<
      ::btool::app::collector::registry::Index>>
      ms_i;
  StrictMock<::btool::app::collector::registry::testing::MockSerializer<
      ::btool::app::collector::registry::Gaggle>>
      ms_g;
  EXPECT_CALL(ms_i, Unmarshal).WillOnce(SetArgPointee<1>(ex_i));
  EXPECT_CALL(ms_g, Unmarshal).WillOnce(SetArgPointee<1>(ex_g0));
  EXPECT_CALL(ms_g, Unmarshal).WillOnce(SetArgPointee<1>(ex_g1));

  ::btool::app::collector::registry::FsRegistry fr(root_, &ms_i, &ms_g);

  ::btool::app::collector::registry::Index ac_i;
  fr.GetIndex(&ac_i);
  EXPECT_EQ(ex_i, ac_i);

  ::btool::app::collector::registry::Gaggle ac_g0;
  fr.GetGaggle("some-path-0", &ac_g0);
  EXPECT_EQ(ex_g0, ac_g0);
  ::btool::app::collector::registry::Gaggle ac_g1;
  fr.GetGaggle("some-path-1", &ac_g1);
  EXPECT_EQ(ex_g1, ac_g1);
}
