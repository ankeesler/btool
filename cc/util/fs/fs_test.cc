#include "fs.h"

#include <string>
#include <vector>

#include "gtest/gtest.h"

#include "core/err.h"

TEST(FS, Base) {
  EXPECT_EQ("", ::btool::util::fs::Base(""));
  EXPECT_EQ("a", ::btool::util::fs::Base("a"));
  EXPECT_EQ(".", ::btool::util::fs::Base("."));
  EXPECT_EQ("b", ::btool::util::fs::Base("a/b"));
  EXPECT_EQ("c.c", ::btool::util::fs::Base("a/b/c.c"));
  EXPECT_EQ("c.c", ::btool::util::fs::Base("./a/b/c.c"));
}

TEST(FS, Dir) {
  EXPECT_EQ("", ::btool::util::fs::Dir(""));
  EXPECT_EQ("a", ::btool::util::fs::Dir("a"));
  EXPECT_EQ(".", ::btool::util::fs::Dir("."));
  EXPECT_EQ("a", ::btool::util::fs::Dir("a/b"));
  EXPECT_EQ("a/b", ::btool::util::fs::Dir("a/b/c.c"));
  EXPECT_EQ("./a/b", ::btool::util::fs::Dir("./a/b/c.c"));
}

TEST(FS, Join) {
  EXPECT_EQ("a/b", ::btool::util::fs::Join("a", "b"));
  EXPECT_EQ("./a/b", ::btool::util::fs::Join("./a", "b"));
  EXPECT_EQ("a/b/c/d/e/f", ::btool::util::fs::Join("a/b/c/d/e", "f"));
}

TEST(FS, Ext) {
  EXPECT_EQ("", ::btool::util::fs::Ext("tuna"));
  EXPECT_EQ("", ::btool::util::fs::Ext("./tuna"));
  EXPECT_EQ("", ::btool::util::fs::Ext("some/path/to/tuna"));

  EXPECT_EQ(".cc", ::btool::util::fs::Ext("tuna.cc"));
  EXPECT_EQ(".cc", ::btool::util::fs::Ext("some/path/to/tuna.cc"));

  EXPECT_EQ(".o", ::btool::util::fs::Ext("tuna.cc.o"));
  EXPECT_EQ(".o", ::btool::util::fs::Ext("some/path/to/tuna.cc.o"));
}

TEST(FS, File) {
  auto err = ::btool::util::fs::TempDir();
  EXPECT_FALSE(err);

  auto dir = err.Ret();
  auto file = ::btool::util::fs::Join(dir, "a");
  EXPECT_TRUE(::btool::util::fs::ReadFile(file));
  EXPECT_TRUE(::btool::util::fs::RemoveAll(file));

  auto void_err =
      ::btool::util::fs::WriteFile(file, "this is text\nwith multiple lines\n");
  EXPECT_FALSE(void_err) << void_err;
  err = ::btool::util::fs::ReadFile(file);
  EXPECT_FALSE(err) << err;
  EXPECT_EQ("this is text\nwith multiple lines\n", err.Ret());

  void_err = ::btool::util::fs::RemoveAll(dir);
  EXPECT_FALSE(void_err) << void_err;
}

TEST(FS, Walk) {
  auto err = ::btool::util::fs::TempDir();
  EXPECT_FALSE(err) << err;
  auto dir = err.Ret();

  std::vector<std::string> paths{
      ::btool::util::fs::Join(dir, "a.c"),
      ::btool::util::fs::Join(dir, "b.h"),

      ::btool::util::fs::Join(dir, "dir0/a.c"),
      ::btool::util::fs::Join(dir, "dir0/b.h"),
      ::btool::util::fs::Join(dir, "dir1/a.c"),
      ::btool::util::fs::Join(dir, "dir1/b.h"),

      ::btool::util::fs::Join(dir, "dir0/dir0/a.c"),
      ::btool::util::fs::Join(dir, "dir0/dir0/b.h"),
      ::btool::util::fs::Join(dir, "dir0/dir1/a.c"),
      ::btool::util::fs::Join(dir, "dir0/dir1/b.h"),
      ::btool::util::fs::Join(dir, "dir1/dir0/a.c"),
      ::btool::util::fs::Join(dir, "dir1/dir0/b.h"),
      ::btool::util::fs::Join(dir, "dir1/dir1/a.c"),
      ::btool::util::fs::Join(dir, "dir1/dir1/b.h"),
  };
  for (auto path : paths) {
    auto dir = ::btool::util::fs::Dir(path);
    auto err = ::btool::util::fs::Exists(dir);
    ASSERT_FALSE(err) << err;
    if (!err.Ret()) {
      auto void_err = ::btool::util::fs::Mkdir(dir);
      ASSERT_FALSE(void_err) << "mkdir " << dir << ": " << void_err;
    }

    auto void_err = ::btool::util::fs::WriteFile(path, "hey\n");
    ASSERT_FALSE(void_err) << void_err;
  }

  std::vector<std::string> ex_visits{
      ::btool::util::fs::Join(dir, "dir0/dir0/a.c"),
      ::btool::util::fs::Join(dir, "dir0/dir0/b.h"),

      ::btool::util::fs::Join(dir, "dir0/dir0"),

      ::btool::util::fs::Join(dir, "dir0/dir1/a.c"),
      ::btool::util::fs::Join(dir, "dir0/dir1/b.h"),

      ::btool::util::fs::Join(dir, "dir0/dir1"),

      ::btool::util::fs::Join(dir, "dir0/a.c"),
      ::btool::util::fs::Join(dir, "dir0/b.h"),

      ::btool::util::fs::Join(dir, "dir0"),

      ::btool::util::fs::Join(dir, "dir1/dir0/a.c"),
      ::btool::util::fs::Join(dir, "dir1/dir0/b.h"),

      ::btool::util::fs::Join(dir, "dir1/dir0"),

      ::btool::util::fs::Join(dir, "dir1/dir1/a.c"),
      ::btool::util::fs::Join(dir, "dir1/dir1/b.h"),

      ::btool::util::fs::Join(dir, "dir1/dir1"),

      ::btool::util::fs::Join(dir, "dir1/a.c"),
      ::btool::util::fs::Join(dir, "dir1/b.h"),

      ::btool::util::fs::Join(dir, "dir1"),

      ::btool::util::fs::Join(dir, "a.c"),
      ::btool::util::fs::Join(dir, "b.h"),

      dir,
  };
  std::vector<std::string> ac_visits;
  auto void_err =
      ::btool::util::fs::Walk(dir, [&ac_visits](const std::string &path) {
        ac_visits.push_back(path);
        return ::btool::core::VoidErr::Success();
      });
  EXPECT_FALSE(void_err) << void_err;
  EXPECT_EQ(ex_visits, ac_visits);

  void_err = ::btool::util::fs::RemoveAll(dir);
  EXPECT_FALSE(void_err) << void_err;
}
