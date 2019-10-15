#include "obj.h"

#include <string>
#include <vector>

// workaround for bug-00
#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "app/collector/store.h"
#include "app/collector/testing/collector.h"
#include "core/err.h"
#include "node/testing/node.h"

using ::testing::Return;

class ObjTest : public ::testing::Test {
 protected:
  ObjTest() : o_(&mrf_) {}

  ::testing::StrictMock<::btool::node::testing::MockResolver> mr_;
  ::testing::StrictMock<::btool::app::collector::testing::MockResolverFactory>
      mrf_;
  ::btool::app::collector::cc::Obj o_;
  ::btool::app::collector::Store s_;
};

TEST_F(ObjTest, Ignore) {
  o_.OnSet(&s_, "foo.go");
  EXPECT_TRUE(s_.IsEmpty());
}

TEST_F(ObjTest, C) {
  std::vector<std::string> include_paths;
  std::vector<std::string> flags;
  EXPECT_CALL(mrf_, NewCompileC(include_paths, flags)).WillOnce(Return(&mr_));

  o_.OnSet(&s_, "foo.c");
  auto n = s_.Get("foo.o");
  EXPECT_TRUE(n);
}

TEST_F(ObjTest, CC) {
  std::vector<std::string> include_paths;
  std::vector<std::string> flags;
  EXPECT_CALL(mrf_, NewCompileCC(include_paths, flags)).WillOnce(Return(&mr_));

  o_.OnSet(&s_, "foo.cc");
  auto n = s_.Get("foo.o");
  EXPECT_TRUE(n);
}
