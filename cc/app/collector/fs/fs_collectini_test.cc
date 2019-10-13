#include "fs_collectini.h"

#include <vector>

#include "gtest/gtest.h"

#include "node/store.h"
#include "util/fs/fs.h"

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
  ::btool::app::collector::fs::FSCollectini fsc(Root());
  ::btool::node::Store s;
  EXPECT_EQ(::btool::core::VoidErr::Success(), fsc.Collect(&s));

  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "a.cc")));
  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "b.h")));
  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "c.c")));
  EXPECT_FALSE(s.Get(::btool::util::fs::Join(Root(), "d.go")));

  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir0/a.cc")));
  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir0/b.h")));
  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir0/c.c")));
  EXPECT_FALSE(s.Get(::btool::util::fs::Join(Root(), "dir0/d.go")));

  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir0/dir0/a.cc")));
  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir0/dir0/b.h")));
  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir0/dir0/c.c")));
  EXPECT_FALSE(s.Get(::btool::util::fs::Join(Root(), "dir0/dir0/d.go")));

  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir0/dir1/a.cc")));
  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir0/dir1/b.h")));
  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir0/dir1/c.c")));
  EXPECT_FALSE(s.Get(::btool::util::fs::Join(Root(), "dir0/dir1/d.go")));

  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir1/a.cc")));
  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir1/b.h")));
  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir1/c.c")));
  EXPECT_FALSE(s.Get(::btool::util::fs::Join(Root(), "dir1/d.go")));

  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir1/dir0/a.cc")));
  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir1/dir0/b.h")));
  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir1/dir0/c.c")));
  EXPECT_FALSE(s.Get(::btool::util::fs::Join(Root(), "dir1/dir0/d.go")));

  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir1/dir1/a.cc")));
  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir1/dir1/b.h")));
  EXPECT_TRUE(s.Get(::btool::util::fs::Join(Root(), "dir1/dir1/c.c")));
  EXPECT_FALSE(s.Get(::btool::util::fs::Join(Root(), "dir1/dir1/d.go")));
}
