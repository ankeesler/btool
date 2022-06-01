#include "app/collector/cc/exe.h"

#include "app/collector/cc/properties.h"
#include "app/collector/cc/testing/cc.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "app/collector/testing/collector.h"
#include "err.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "node/testing/node.h"

using ::testing::ElementsAre;
using ::testing::Return;

class ExeTest : public ::testing::Test {
 protected:
  ExeTest() : e_(&mrf_) {}

  ::testing::StrictMock<::btool::node::testing::MockResolver> mr_;
  ::testing::StrictMock<
      ::btool::app::collector::cc::testing::MockResolverFactory>
      mrf_;
  ::btool::app::collector::cc::Exe e_;
  ::btool::app::collector::Store s_;
};

void BuildGraph(::btool::app::collector::Store *s, const std::string ext);

TEST_F(ExeTest, BadFileExtension) {
  s_.Put("tuna.go");
  e_.OnNotify(&s_, "tuna.go");
  EXPECT_EQ(1UL, s_.Size());
}

TEST_F(ExeTest, C) {
  EXPECT_CALL(mrf_, NewLinkC()).WillOnce(Return(&mr_));

  BuildGraph(&s_, ".c");

  ::btool::app::collector::testing::SpyCollectini sc;

  e_.OnNotify(&s_, "tuna");
  EXPECT_EQ(16UL, s_.Size());

  auto n = s_.Get("tuna");
  EXPECT_EQ(6UL, n->dependencies()->size());
  EXPECT_EQ("tuna.o", n->dependencies()->at(0)->name());
  EXPECT_EQ("c.o", n->dependencies()->at(1)->name());
  EXPECT_EQ("a.o", n->dependencies()->at(2)->name());
  EXPECT_EQ("marlin.o", n->dependencies()->at(3)->name());
  EXPECT_EQ("b.o", n->dependencies()->at(4)->name());
  EXPECT_EQ("lib.a", n->dependencies()->at(5)->name());
  EXPECT_EQ(&mr_, n->resolver());

  EXPECT_THAT(
      sc.on_notify_calls_,
      ElementsAre(
          std::pair<::btool::app::collector::Store *, const std::string &>(
              &s_, n->name())));
}

TEST_F(ExeTest, CC) {
  EXPECT_CALL(mrf_, NewLinkCC()).WillOnce(Return(&mr_));

  BuildGraph(&s_, ".cc");

  ::btool::app::collector::testing::SpyCollectini sc;

  e_.OnNotify(&s_, "tuna");
  EXPECT_EQ(16UL, s_.Size());

  auto n = s_.Get("tuna");
  EXPECT_EQ(6UL, n->dependencies()->size());
  EXPECT_EQ("tuna.o", n->dependencies()->at(0)->name());
  EXPECT_EQ("c.o", n->dependencies()->at(1)->name());
  EXPECT_EQ("a.o", n->dependencies()->at(2)->name());
  EXPECT_EQ("marlin.o", n->dependencies()->at(3)->name());
  EXPECT_EQ("b.o", n->dependencies()->at(4)->name());
  EXPECT_EQ("lib.a", n->dependencies()->at(5)->name());
  EXPECT_EQ(&mr_, n->resolver());

  EXPECT_THAT(
      sc.on_notify_calls_,
      ElementsAre(
          std::pair<::btool::app::collector::Store *, const std::string &>(
              &s_, n->name())));
}

void BuildGraph(::btool::app::collector::Store *s, const std::string ext) {
  // tuna -> a -> b -> marlin
  // tuna -> marlin
  // a -> c

  s->Put("lib.a");

  auto marlinh = s->Put("marlin.h");
  ::btool::app::collector::cc::Properties::AddLibrary(marlinh->property_store(),
                                                      "lib.a");
  auto marlinc = s->Put("marlin" + ext);
  marlinc->dependencies()->push_back(marlinh);
  auto marlino = s->Put("marlin.o");
  marlino->dependencies()->push_back(marlinc);

  auto ch = s->Put("c.h");
  auto cc = s->Put("c" + ext);
  cc->dependencies()->push_back(ch);
  auto co = s->Put("c.o");
  co->dependencies()->push_back(cc);

  auto bh = s->Put("b.h");
  bh->dependencies()->push_back(marlinh);
  auto bc = s->Put("b" + ext);
  bc->dependencies()->push_back(bh);
  auto bo = s->Put("b.o");
  bo->dependencies()->push_back(bc);

  auto ah = s->Put("a.h");
  ah->dependencies()->push_back(ch);
  auto ac = s->Put("a" + ext);
  ac->dependencies()->push_back(ah);
  ac->dependencies()->push_back(bh);
  auto ao = s->Put("a.o");
  ao->dependencies()->push_back(ac);

  auto tunac = s->Put("tuna" + ext);
  tunac->dependencies()->push_back(ah);
  tunac->dependencies()->push_back(marlinh);
  auto tunao = s->Put("tuna.o");
  tunao->dependencies()->push_back(tunac);

  s->Put("tuna");
}
