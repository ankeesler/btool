#include "obj.h"

#include <string>
#include <vector>

// workaround for bug-00
#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "app/collector/cc/properties.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "app/collector/testing/collector.h"
#include "core/err.h"
#include "node/testing/node.h"

using ::testing::Return;

class ObjTest : public ::testing::Test {
 protected:
  ObjTest() : o_(&mrf_) {
    auto bh = s_.Put("b.h");
    ::btool::app::collector::cc::Properties::AddIncludePath(
        bh->property_store(), "bh/include/path");

    auto ah = s_.Put("a.h");
    ah->dependencies()->push_back(bh);
    ::btool::app::collector::cc::Properties::AddIncludePath(
        ah->property_store(), "ah/include/path");

    auto fooc = s_.Put("foo.c");
    fooc->dependencies()->push_back(ah);
    ::btool::app::collector::Properties::SetLocal(fooc->property_store(), true);
    ::btool::app::collector::cc::Properties::AddIncludePath(
        fooc->property_store(), "fooc/include/path");

    auto foocc = s_.Put("foo.cc");
    foocc->dependencies()->push_back(ah);
    ::btool::app::collector::Properties::SetLocal(foocc->property_store(),
                                                  true);
    ::btool::app::collector::cc::Properties::AddIncludePath(
        foocc->property_store(), "foocc/include/path");

    auto foogo = s_.Put("foo.go");
    ::btool::app::collector::Properties::SetLocal(foogo->property_store(),
                                                  true);

    s_.Put("bar.c");
  }

  ::testing::StrictMock<::btool::node::testing::MockResolver> mr_;
  ::testing::StrictMock<::btool::app::collector::testing::MockResolverFactory>
      mrf_;
  ::btool::app::collector::cc::Obj o_;
  ::btool::app::collector::Store s_;
};

TEST_F(ObjTest, IgnoreFileExt) {
  o_.OnSet(&s_, "foo.go");
  EXPECT_EQ(6, s_.Size());
}

TEST_F(ObjTest, IgnoreNotLocal) {
  o_.OnSet(&s_, "bar.c");
  EXPECT_EQ(6, s_.Size());
}

TEST_F(ObjTest, C) {
  std::vector<std::string> include_paths{"bh/include/path", "ah/include/path",
                                         "fooc/include/path"};
  std::vector<std::string> flags;
  EXPECT_CALL(mrf_, NewCompileC(include_paths, flags)).WillOnce(Return(&mr_));

  auto d = s_.Get("foo.c");
  o_.OnSet(&s_, d->name());

  auto n = s_.Get("foo.o");
  std::vector<::btool::node::Node *> ex_deps{d};
  EXPECT_TRUE(n);
  EXPECT_EQ(ex_deps, *n->dependencies());
  EXPECT_EQ(&mr_, n->resolver());
}

TEST_F(ObjTest, CC) {
  std::vector<std::string> include_paths{"bh/include/path", "ah/include/path",
                                         "foocc/include/path"};
  std::vector<std::string> flags;
  EXPECT_CALL(mrf_, NewCompileCC(include_paths, flags)).WillOnce(Return(&mr_));

  auto d = s_.Get("foo.cc");
  o_.OnSet(&s_, d->name());

  auto n = s_.Get("foo.o");
  std::vector<::btool::node::Node *> ex_deps{d};
  EXPECT_TRUE(n);
  EXPECT_EQ(ex_deps, *n->dependencies());
  EXPECT_EQ(&mr_, n->resolver());
}
