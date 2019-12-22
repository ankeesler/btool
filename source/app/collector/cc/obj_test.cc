#include "obj.h"

#include <string>
#include <vector>

#include "gtest/gtest.h"

#include "app/collector/cc/properties.h"
#include "app/collector/cc/testing/cc.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "app/collector/testing/collector.h"
#include "err.h"
#include "node/testing/node.h"

using ::testing::ElementsAre;
using ::testing::Return;

class ObjTest : public ::testing::Test {
 protected:
  ObjTest() : o_(&mrf_) {
    auto bh = s_.Put("b.h");

    auto ah = s_.Put("a.h");
    ah->dependencies()->push_back(bh);

    auto fooc = s_.Put("foo.c");
    fooc->dependencies()->push_back(ah);
    ::btool::app::collector::Properties::SetLocal(fooc->property_store(), true);

    auto foocc = s_.Put("foo.cc");
    foocc->dependencies()->push_back(ah);
    ::btool::app::collector::Properties::SetLocal(foocc->property_store(),
                                                  true);

    auto foogo = s_.Put("foo.go");
    ::btool::app::collector::Properties::SetLocal(foogo->property_store(),
                                                  true);

    s_.Put("bar.c");
  }

  ::testing::StrictMock<::btool::node::testing::MockResolver> mr_;
  ::testing::StrictMock<
      ::btool::app::collector::cc::testing::MockResolverFactory>
      mrf_;
  ::btool::app::collector::cc::Obj o_;
  ::btool::app::collector::Store s_;
};

TEST_F(ObjTest, IgnoreFileExt) {
  o_.OnNotify(&s_, "foo.go");
  EXPECT_EQ(6UL, s_.Size());
}

TEST_F(ObjTest, IgnoreNotLocal) {
  o_.OnNotify(&s_, "bar.c");
  EXPECT_EQ(6UL, s_.Size());
}

TEST_F(ObjTest, C) {
  EXPECT_CALL(mrf_, NewCompileC()).WillOnce(Return(&mr_));

  ::btool::app::collector::testing::SpyCollectini sc;

  auto d = s_.Get("foo.c");
  o_.OnNotify(&s_, d->name());

  auto n = s_.Get("foo.o");
  std::vector<::btool::node::Node *> ex_deps{d};
  EXPECT_TRUE(n);
  EXPECT_EQ(ex_deps, *n->dependencies());
  EXPECT_EQ(&mr_, n->resolver());

  EXPECT_THAT(
      sc.on_notify_calls_,
      ElementsAre(
          std::pair<::btool::app::collector::Store *, const std::string &>(
              &s_, n->name())));
}

TEST_F(ObjTest, CC) {
  EXPECT_CALL(mrf_, NewCompileCC()).WillOnce(Return(&mr_));

  ::btool::app::collector::testing::SpyCollectini sc;

  auto d = s_.Get("foo.cc");
  o_.OnNotify(&s_, d->name());

  auto n = s_.Get("foo.o");
  std::vector<::btool::node::Node *> ex_deps{d};
  EXPECT_TRUE(n);
  EXPECT_EQ(ex_deps, *n->dependencies());
  EXPECT_EQ(&mr_, n->resolver());

  EXPECT_THAT(
      sc.on_notify_calls_,
      ElementsAre(
          std::pair<::btool::app::collector::Store *, const std::string &>(
              &s_, n->name())));
}
