#include "app/collector/cc/exe.h"

#include "gmock/gmock.h"
#include "gtest/gtest.h"

// workaround for bug-02
#include "app/collector/base_collectini.h"
#include "app/collector/cc/properties.h"
#include "app/collector/cc/testing/cc.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "app/collector/testing/collector.h"
#include "core/err.h"
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
  std::vector<std::string> flags;
  flags.push_back("-some-flag");
  EXPECT_CALL(mrf_, NewLinkC(flags)).WillOnce(Return(&mr_));

  BuildGraph(&s_, ".c");

  ::btool::app::collector::testing::SpyCollectini sc;

  e_.OnNotify(&s_, "tuna");
  EXPECT_EQ(13UL, s_.Size());

  auto n = s_.Get("tuna");
  EXPECT_EQ(5UL, n->dependencies()->size());
  EXPECT_EQ("tuna.o", n->dependencies()->at(0)->name());
  EXPECT_EQ("a.o", n->dependencies()->at(1)->name());
  EXPECT_EQ("b.o", n->dependencies()->at(2)->name());
  EXPECT_EQ("marlin.o", n->dependencies()->at(3)->name());
  EXPECT_EQ("lib.a", n->dependencies()->at(4)->name());
  EXPECT_EQ(&mr_, n->resolver());

  EXPECT_THAT(
      sc.on_notify_calls_,
      ElementsAre(
          std::pair<::btool::app::collector::Store *, const std::string &>(
              &s_, n->name())));
}

TEST_F(ExeTest, CC) {
  std::vector<std::string> flags;
  flags.push_back("-some-flag");
  EXPECT_CALL(mrf_, NewLinkCC(flags)).WillOnce(Return(&mr_));

  BuildGraph(&s_, ".cc");

  ::btool::app::collector::testing::SpyCollectini sc;

  e_.OnNotify(&s_, "tuna");
  EXPECT_EQ(13UL, s_.Size());

  auto n = s_.Get("tuna");
  EXPECT_EQ(5UL, n->dependencies()->size());
  EXPECT_EQ("tuna.o", n->dependencies()->at(0)->name());
  EXPECT_EQ("a.o", n->dependencies()->at(1)->name());
  EXPECT_EQ("b.o", n->dependencies()->at(2)->name());
  EXPECT_EQ("marlin.o", n->dependencies()->at(3)->name());
  EXPECT_EQ("lib.a", n->dependencies()->at(4)->name());
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

  auto liba = s->Put("lib.a");
  ::btool::app::collector::cc::Properties::AddLinkFlag(liba->property_store(),
                                                       "-some-flag");

  auto marlinh = s->Put("marlin.h");
  ::btool::app::collector::cc::Properties::AddLibrary(marlinh->property_store(),
                                                      "lib.a");
  auto marlinc = s->Put("marlin" + ext);
  marlinc->dependencies()->push_back(marlinh);
  auto marlino = s->Put("marlin.o");
  marlino->dependencies()->push_back(marlinc);

  auto bh = s->Put("b.h");
  bh->dependencies()->push_back(marlinh);
  auto bc = s->Put("b" + ext);
  bc->dependencies()->push_back(bh);
  auto bo = s->Put("b.o");
  bo->dependencies()->push_back(bc);

  auto ah = s->Put("a.h");
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
