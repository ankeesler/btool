#include "app/collector/cc/inc.h"

#include <functional>
#include <string>
#include <vector>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "app/collector/cc/properties.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "core/err.h"

using ::testing::_;
using ::testing::Return;
using ::testing::StrictMock;

class MockIncludesParser
    : public ::btool::app::collector::cc::Inc::IncludesParser {
 public:
  MOCK_METHOD1(ParseIncludes, ::btool::core::VoidErr(
                                  std::function<void(const std::string &)>));
};

class FakeIncludesParser
    : public ::btool::app::collector::cc::Inc::IncludesParser {
 public:
  void AddInclude(const std::string &include) { includes_.push_back(include); }

  ::btool::core::VoidErr ParseIncludes(
      std::function<void(const std::string &)> f) override {
    for (auto include : includes_) {
      f(include);
    }
    return ::btool::core::VoidErr::Success();
  }

 private:
  std::vector<std::string> includes_;
};

TEST(IncTest, NotLocal) {
  ::btool::app::collector::Store s;
  s.Put("tuna.h");

  StrictMock<MockIncludesParser> mip;
  ::btool::app::collector::cc::Inc i(&mip);
  i.OnSet(&s, "tuna.h");

  EXPECT_EQ(0UL, s.Get("tuna.h")->dependencies()->size());
}

TEST(IncTest, BadFileExt) {
  ::btool::app::collector::Store s;
  auto n = s.Put("tuna.go");
  ::btool::app::collector::Properties::SetLocal(n->property_store(), true);

  StrictMock<MockIncludesParser> mip;
  ::btool::app::collector::cc::Inc i(&mip);
  i.OnSet(&s, "tuna.go");

  EXPECT_EQ(0UL, s.Get("tuna.go")->dependencies()->size());
}

TEST(IncTest, C) {
  ::btool::app::collector::Store s;
  auto n = s.Put("tuna.c");
  ::btool::app::collector::Properties::SetLocal(n->property_store(), true);
  s.Put("some/root/some/path.h");
  s.Put("some/root/some/other/path.h");
  s.Put("some/lib/include/lib/path.h");

  FakeIncludesParser fip;
  fip.AddInclude("some/path.h");
  fip.AddInclude("some/other/path.h");
  fip.AddInclude("lib/path.h");
  ::btool::app::collector::cc::Inc i(&fip);
  i.OnSet(&s, "tuna.c");

  auto deps = s.Get("tuna.c")->dependencies();
  EXPECT_EQ(3UL, deps->size());
  EXPECT_EQ("some/root/some/path.h", deps->at(0)->name());
  EXPECT_EQ("some/root/some/other/path.h", deps->at(1)->name());
  EXPECT_EQ("some/lib/include/lib/path.h", deps->at(2)->name());

  auto include_paths = ::btool::app::collector::cc::Properties::IncludePaths(
      s.Get("tuna.c")->property_store());
  EXPECT_EQ(3UL, include_paths->size());
  EXPECT_EQ("some/root/", include_paths->at(0));
  EXPECT_EQ("some/root/", include_paths->at(1));
  EXPECT_EQ("some/lib/include/", include_paths->at(2));
}

TEST(IncTest, EmptyIncludePath) {
  ::btool::app::collector::Store s;
  auto n = s.Put("tuna.c");
  ::btool::app::collector::Properties::SetLocal(n->property_store(), true);
  s.Put("some/path.h");

  FakeIncludesParser fip;
  fip.AddInclude("some/path.h");
  ::btool::app::collector::cc::Inc i(&fip);
  i.OnSet(&s, "tuna.c");

  auto deps = s.Get("tuna.c")->dependencies();
  EXPECT_EQ(1UL, deps->size());
  EXPECT_EQ("some/path.h", deps->at(0)->name());

  auto include_paths = ::btool::app::collector::cc::Properties::IncludePaths(
      s.Get("tuna.c")->property_store());
  EXPECT_EQ(1UL, include_paths->size());
  EXPECT_EQ(".", include_paths->at(0));
}
