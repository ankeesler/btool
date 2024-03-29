#include "app/collector/cc/inc.h"

#include <functional>
#include <string>
#include <utility>
#include <vector>

#include "app/collector/cc/properties.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "app/collector/testing/collector.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"

using ::testing::_;
using ::testing::ElementsAre;
using ::testing::Return;
using ::testing::StrictMock;

class MockIncludesParser
    : public ::btool::app::collector::cc::Inc::IncludesParser {
 public:
  MOCK_METHOD2(ParseIncludes,
               void(const std::string &,
                    std::function<void(const std::string &, bool)>));
};

class FakeIncludesParser
    : public ::btool::app::collector::cc::Inc::IncludesParser {
 public:
  void AddInclude(std::string include, bool system) {
    includes_.push_back(std::pair<std::string, bool>(include, system));
  }

  void ParseIncludes(
      const std::string &name,
      std::function<void(const std::string &, bool)> f) override {
    for (const auto &include : includes_) {
      f(include.first, include.second);
    }
  }

 private:
  std::vector<std::pair<std::string, bool>> includes_;
};

TEST(IncTest, NotLocal) {
  ::btool::app::collector::Store s;
  s.Put("tuna.h");

  StrictMock<MockIncludesParser> mip;
  ::btool::app::collector::cc::Inc i(&mip);
  i.OnNotify(&s, "tuna.h");

  EXPECT_EQ(0UL, s.Get("tuna.h")->dependencies()->size());
}

TEST(IncTest, BadFileExt) {
  ::btool::app::collector::Store s;
  auto n = s.Put("tuna.go");
  ::btool::app::collector::Properties::SetLocal(n->property_store(), true);

  StrictMock<MockIncludesParser> mip;
  ::btool::app::collector::cc::Inc i(&mip);
  i.OnNotify(&s, "tuna.go");

  EXPECT_EQ(0UL, s.Get("tuna.go")->dependencies()->size());
}

TEST(IncTest, C) {
  ::btool::app::collector::Store s;
  auto n = s.Put("tuna.c");
  ::btool::app::collector::Properties::SetLocal(n->property_store(), true);
  s.Put("some/root/some/path.h");
  s.Put("some/root/some/other/path.h");
  s.Put("some/lib/include/lib/path.h");

  ::btool::app::collector::testing::SpyCollectini sc;
  FakeIncludesParser fip;
  fip.AddInclude("string", true);
  fip.AddInclude("some/path.h", false);
  fip.AddInclude("some/other/path.h", false);
  fip.AddInclude("lib/path.h", false);
  ::btool::app::collector::cc::Inc i(&fip);
  i.OnNotify(&s, "tuna.c");

  auto deps = s.Get("tuna.c")->dependencies();
  EXPECT_EQ(3UL, deps->size());
  EXPECT_EQ("some/root/some/path.h", deps->at(0)->name());
  EXPECT_EQ("some/root/some/other/path.h", deps->at(1)->name());
  EXPECT_EQ("some/lib/include/lib/path.h", deps->at(2)->name());

  auto include_paths = ::btool::app::collector::cc::Properties::IncludePaths(
      s.Get("tuna.c")->property_store());
  EXPECT_EQ(2UL, include_paths->size());
  EXPECT_EQ("some/root/", include_paths->at(0));
  EXPECT_EQ("some/lib/include/", include_paths->at(1));

  EXPECT_THAT(
      sc.on_notify_calls_,
      ElementsAre(
          std::pair<::btool::app::collector::Store *, const std::string &>(
              &s, n->name())));

  i.OnNotify(&s, "tuna.c");

  deps = s.Get("tuna.c")->dependencies();
  EXPECT_EQ(3UL, deps->size());
  EXPECT_EQ("some/root/some/path.h", deps->at(0)->name());
  EXPECT_EQ("some/root/some/other/path.h", deps->at(1)->name());
  EXPECT_EQ("some/lib/include/lib/path.h", deps->at(2)->name());

  include_paths = ::btool::app::collector::cc::Properties::IncludePaths(
      s.Get("tuna.c")->property_store());
  EXPECT_EQ(2UL, include_paths->size());
  EXPECT_EQ("some/root/", include_paths->at(0));
  EXPECT_EQ("some/lib/include/", include_paths->at(1));

  EXPECT_THAT(
      sc.on_notify_calls_,
      ElementsAre(
          std::pair<::btool::app::collector::Store *, const std::string &>(
              &s, n->name())));
}

TEST(IncTest, EmptyIncludePath) {
  ::btool::app::collector::Store s;
  auto n = s.Put("tuna.c");
  ::btool::app::collector::Properties::SetLocal(n->property_store(), true);
  s.Put("some/path.h");

  ::btool::app::collector::testing::SpyCollectini sc;
  FakeIncludesParser fip;
  fip.AddInclude("some/path.h", false);
  ::btool::app::collector::cc::Inc i(&fip);
  i.OnNotify(&s, "tuna.c");

  auto deps = s.Get("tuna.c")->dependencies();
  EXPECT_EQ(1UL, deps->size());
  EXPECT_EQ("some/path.h", deps->at(0)->name());

  auto include_paths = ::btool::app::collector::cc::Properties::IncludePaths(
      s.Get("tuna.c")->property_store());
  EXPECT_EQ(1UL, include_paths->size());
  EXPECT_EQ(".", include_paths->at(0));

  EXPECT_THAT(
      sc.on_notify_calls_,
      ElementsAre(
          std::pair<::btool::app::collector::Store *, const std::string &>(
              &s, n->name())));
}

TEST(IncTest, Pthread) {
  ::btool::app::collector::Store s;
  auto n = s.Put("tuna.c");
  ::btool::app::collector::Properties::SetLocal(n->property_store(), true);

  ::btool::app::collector::testing::SpyCollectini sc;
  FakeIncludesParser fip;
  fip.AddInclude("thread", true);
  ::btool::app::collector::cc::Inc i(&fip);
  i.OnNotify(&s, "tuna.c");

  auto deps = s.Get("tuna.c")->dependencies();
  EXPECT_EQ(0UL, deps->size());

  auto include_paths = ::btool::app::collector::cc::Properties::IncludePaths(
      s.Get("tuna.c")->property_store());
  EXPECT_EQ(nullptr, include_paths);

  auto link_flags = ::btool::app::collector::cc::Properties::LinkFlags(
      s.Get("tuna.c")->property_store());
  EXPECT_THAT(*link_flags, ElementsAre("-lpthread"));

  EXPECT_THAT(
      sc.on_notify_calls_,
      ElementsAre(
          std::pair<::btool::app::collector::Store *, const std::string &>(
              &s, n->name())));
}
