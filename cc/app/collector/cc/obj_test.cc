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

TEST_F(ObjTest, IgnoreFileExt) {
  auto n = s_.Put("foo.go");
  o_.OnSet(&s_, "foo.go");
  EXPECT_TRUE(s_.IsEmpty());
}

TEST_F(ObjTest, IgnoreNotLocal) {
  o_.OnSet(&s_, "foo.cc");
  EXPECT_TRUE(s_.IsEmpty());
}

TEST_F(ObjTest, C) {
  std::vector<std::string> include_paths;
  std::vector<std::string> flags;
  EXPECT_CALL(mrf_, NewCompileC(include_paths, flags)).WillOnce(Return(&mr_));

  auto d = s_.Put("foo.c");
  o_.OnSet(&s_, d->Name());

  auto n = s_.Get("foo.o");
  std::vector<::btool::node::Node *> ex_deps{d};
  EXPECT_TRUE(n);
  EXPECT_EQ(ex_deps, n->Deps());
  EXPECT_EQ(&mr_, n->Resolver());
}

TEST_F(ObjTest, CC) {
  std::vector<std::string> include_paths;
  std::vector<std::string> flags;
  EXPECT_CALL(mrf_, NewCompileCC(include_paths, flags)).WillOnce(Return(&mr_));

  auto d = s_.Put("foo.cc");
  o_.OnSet(&s_, d->Name());

  auto n = s_.Get("foo.o");
  std::vector<::btool::node::Node *> ex_deps{d};
  EXPECT_TRUE(n);
  EXPECT_EQ(ex_deps, n->Deps());
  EXPECT_EQ(&mr_, n->Resolver());
}
