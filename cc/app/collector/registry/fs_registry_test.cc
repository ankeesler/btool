#include "app/collector/registry/fs_registry.h"

#include <string>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "app/collector/registry/testing/registry.h"
#include "util/fs/fs.h"

using ::testing::ContainerEq;
using ::testing::SetArgPointee;
using ::testing::StrictMock;

class FsRegistryTest : public ::testing::Test {
 protected:
  void SetUp() override {
    root_ = ::btool::util::fs::TempDir();

    auto index_yml = ::btool::util::fs::Join(root_, "index.yml");
    ::btool::util::fs::WriteFile(index_yml, "here is an index.yml file\n");
  }

  void TearDown() override { ::btool::util::fs::RemoveAll(root_); }

  std::string root_;
};

TEST_F(FsRegistryTest, GetIndex) {
  ::btool::app::collector::registry::IndexFile if0{.path = "some-path-0",
                                                   .sha256 = "sha0"};
  ::btool::app::collector::registry::IndexFile if1{.path = "some-path-1",
                                                   .sha256 = "sha1"};
  ::btool::app::collector::registry::Index ex_i{.files = {if0, if1}};

  StrictMock<::btool::app::collector::registry::testing::MockSerializer> ms;
  EXPECT_CALL(ms, UnmarshalIndex).WillOnce(SetArgPointee<1>(ex_i));

  ::btool::app::collector::registry::FsRegistry fr(root_, &ms);
  ::btool::app::collector::registry::Index ac_i;
  fr.GetIndex(&ac_i);
  EXPECT_EQ(ex_i, ac_i);
}
