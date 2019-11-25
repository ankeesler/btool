#include "fs.h"

#include <string>
#include <vector>

#include "gtest/gtest.h"

#include "err.h"

TEST(FS, Base) {
  EXPECT_EQ("", ::btool::util::fs::Base(""));
  EXPECT_EQ("a", ::btool::util::fs::Base("a"));
  EXPECT_EQ(".", ::btool::util::fs::Base("."));
  EXPECT_EQ("b", ::btool::util::fs::Base("a/b"));
  EXPECT_EQ("c.c", ::btool::util::fs::Base("a/b/c.c"));
  EXPECT_EQ("c.c", ::btool::util::fs::Base("./a/b/c.c"));
}

TEST(FS, Dir) {
  EXPECT_EQ(".", ::btool::util::fs::Dir(""));
  EXPECT_EQ(".", ::btool::util::fs::Dir("a"));
  EXPECT_EQ(".", ::btool::util::fs::Dir("."));
  EXPECT_EQ("a", ::btool::util::fs::Dir("a/b"));
  EXPECT_EQ("a/b", ::btool::util::fs::Dir("a/b/c.c"));
  EXPECT_EQ("./a/b", ::btool::util::fs::Dir("./a/b/c.c"));

  EXPECT_EQ(".", ::btool::util::fs::Dir("c.c"));
  EXPECT_EQ(".", ::btool::util::fs::Dir("./c.c"));

  EXPECT_EQ("/", ::btool::util::fs::Dir("/tmp"));
  EXPECT_EQ("/", ::btool::util::fs::Dir("/"));
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
  auto dir = ::btool::util::fs::TempDir();

  auto file = ::btool::util::fs::Join(dir, "a");
  EXPECT_THROW(::btool::util::fs::ReadFile(file), ::btool::Err);
  ::btool::util::fs::RemoveAll(file);

  ::btool::util::fs::WriteFile(file, "this is text\nwith multiple lines\n");
  auto content = ::btool::util::fs::ReadFile(file);
  EXPECT_EQ("this is text\nwith multiple lines\n", content);

  ::btool::util::fs::RemoveAll(dir);
}

TEST(FS, MkdirAll) {
  auto dir = ::btool::util::fs::TempDir();
  auto a = ::btool::util::fs::Join(dir, "a");
  auto b = ::btool::util::fs::Join(a, "b");
  auto c = ::btool::util::fs::Join(b, "c");
  ::btool::util::fs::MkdirAll(c);
  EXPECT_TRUE(::btool::util::fs::Exists(c));
  ::btool::util::fs::RemoveAll(dir);
}

TEST(FS, Walk) {
  auto dir = ::btool::util::fs::TempDir();

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
  for (const auto &path : paths) {
    auto dir = ::btool::util::fs::Dir(path);
    auto exists = ::btool::util::fs::Exists(dir);
    if (!exists) {
      ::btool::util::fs::Mkdir(dir);
    }

    ::btool::util::fs::WriteFile(path, "hey\n");
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
  ::btool::util::fs::Walk(dir, [&ac_visits](const std::string &path) {
    ac_visits.push_back(path);
  });
  EXPECT_EQ(ex_visits, ac_visits);

  ::btool::util::fs::RemoveAll(dir);
}

TEST(FS, Is) {
  auto dir = ::btool::util::fs::TempDir();
  auto file = ::btool::util::fs::Join(dir, "a");
  ::btool::util::fs::WriteFile(file, "hey\n");

  EXPECT_TRUE(::btool::util::fs::IsDir(dir));
  // EXPECT_FALSE(::btool::util::fs::IsFile(dir));

  EXPECT_FALSE(::btool::util::fs::IsDir(file));
  // EXPECT_TRUE(::btool::util::fs::IsFile(file));

  ::btool::util::fs::RemoveAll(dir);
}
