#include "fs_collectini.h"

#include <vector>

// workaround for bug-00
#include "gmock/gmock.h"
#include "gtest/gtest.h"

// workaround for bug-02
#include "app/collector/base_collectini.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "app/collector/testing/collector.h"
#include "util/fs/fs.h"

using ::testing::Contains;

class FSCollectiniTest : public ::testing::Test {
 protected:
  void SetUp() override {
    auto err = ::btool::util::fs::TempDir();
    ASSERT_FALSE(err) << err;
    root_ = err.Ret();

    const std::vector<std::string> dirs{
        root_,
        ::btool::util::fs::Join(root_, "dir0"),
        ::btool::util::fs::Join(root_, "dir1"),
        ::btool::util::fs::Join(root_, "dir0/dir0"),
        ::btool::util::fs::Join(root_, "dir0/dir1"),
        ::btool::util::fs::Join(root_, "dir1/dir0"),
        ::btool::util::fs::Join(root_, "dir1/dir1"),
    };
    const std::vector<std::string> files{
        "a.cc",
        "b.h",
        "c.c",
        "d.go",
    };
    for (auto dir : dirs) {
      for (auto file : files) {
        auto err = ::btool::util::fs::Exists(dir);
        ASSERT_FALSE(err) << err;
        if (!err.Ret()) {
          auto void_err = ::btool::util::fs::Mkdir(dir);
          ASSERT_FALSE(void_err) << void_err;
        }

        auto path = ::btool::util::fs::Join(dir, file);
        auto void_err = ::btool::util::fs::WriteFile(path, "hey\n");
        ASSERT_FALSE(void_err) << void_err;
      }
    }
  }

  void TearDown() override {
    auto err = ::btool::util::fs::RemoveAll(root_);
    ASSERT_FALSE(err) << err;
  }

  const std::string &Root() const { return root_; };

 private:
  std::string root_;
};

TEST_F(FSCollectiniTest, Yeah) {
  ::btool::app::collector::testing::SpyCollectini sc;
  ::btool::app::collector::fs::FSCollectini fsc(Root());
  ::btool::app::collector::Store s;
  fsc.Collect(&s);

  const std::vector<std::string> yes{
      "a.cc",           "b.h",           "c.c",

      "dir0/a.cc",      "dir0/b.h",      "dir0/c.c",

      "dir0/dir0/a.cc", "dir0/dir0/b.h", "dir0/dir0/c.c",

      "dir0/dir1/a.cc", "dir0/dir1/b.h", "dir0/dir1/c.c",

      "dir1/a.cc",      "dir1/b.h",      "dir1/c.c",

      "dir1/dir0/a.cc", "dir1/dir0/b.h", "dir1/dir0/c.c",

      "dir1/dir1/a.cc", "dir1/dir1/b.h", "dir1/dir1/c.c",
  };

  const std::vector<std::string> no{
      "d.go",      "dir0/d.go",      "dir0/dir0/d.go", "dir0/dir1/d.go",
      "dir1/d.go", "dir1/dir0/d.go", "dir1/dir1/d.go",

  };

  for (auto f : yes) {
    auto n = s.Get(::btool::util::fs::Join(Root(), f));
    EXPECT_TRUE(n != nullptr);
    EXPECT_TRUE(
        ::btool::app::collector::Properties::Local(n->property_store()));
    EXPECT_THAT(
        sc.on_notify_calls_,
        Contains(
            std::pair<::btool::app::collector::Store *, const std::string &>(
                &s, n->name())));
  }

  for (auto f : no) {
    auto n = s.Get(::btool::util::fs::Join(Root(), f));
    EXPECT_TRUE(n == nullptr);
  }
}
