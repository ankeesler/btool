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
  std::string dir;
  std::string err;
  ASSERT_TRUE(::btool::util::fs::TempDir(&dir, &err)) << "error: " << err;

  auto file = ::btool::util::fs::Join(dir, "a");
  std::string content;
  EXPECT_FALSE(::btool::util::fs::ReadFile(file, &content, &err));
  EXPECT_TRUE(::btool::util::fs::RemoveAll(file, &err)) << "error: " << err;

  EXPECT_TRUE(::btool::util::fs::WriteFile(
      file, "this is text\nwith multiple lines\n", &err))
      << "error: " << err;
  EXPECT_TRUE(::btool::util::fs::ReadFile(file, &content, &err))
      << "error: " << err;
  EXPECT_EQ("this is text\nwith multiple lines\n", content);

  EXPECT_TRUE(::btool::util::fs::RemoveAll(dir, &err)) << "error: " << err;
}

TEST(FS, Walk) {
  std::string dir;
  std::string err;
  ASSERT_TRUE(::btool::util::fs::TempDir(&dir, &err)) << "error: " << err;

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
    bool exists;
    ASSERT_TRUE(::btool::util::fs::Exists(dir, &exists, &err))
        << "error: " << err;
    if (!exists) {
      ASSERT_TRUE(::btool::util::fs::Mkdir(dir, &err)) << "error: " << err;
    }

    ASSERT_TRUE(::btool::util::fs::WriteFile(path, "hey\n", &err))
        << "error: " << err;
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
  EXPECT_TRUE(::btool::util::fs::Walk(
      dir,
      [&ac_visits](const std::string &path, std::string *ret_err) -> bool {
        ac_visits.push_back(path);
        return true;
      },
      &err))
      << "error: " << err;
  EXPECT_EQ(ex_visits, ac_visits);

  EXPECT_TRUE(::btool::util::fs::RemoveAll(dir, &err)) << "error: " << err;
}

TEST(FS, Is) {
  std::string dir;
  std::string err;
  ASSERT_TRUE(::btool::util::fs::TempDir(&dir, &err)) << "error: " << err;

  auto file = ::btool::util::fs::Join(dir, "a");
  ASSERT_TRUE(::btool::util::fs::WriteFile(file, "hey\n", &err))
      << "error: " << err;

  bool is_dir;
  EXPECT_TRUE(::btool::util::fs::IsDir(dir, &is_dir, &err)) << "error: " << err;
  EXPECT_TRUE(is_dir);
  //  EXPECT_FALSE(::btool::util::fs::IsFile(dir));

  EXPECT_TRUE(::btool::util::fs::IsDir(file, &is_dir, &err))
      << "error: " << err;
  EXPECT_FALSE(is_dir);
  // EXPECT_TRUE(::btool::util::fs::IsFile(file));

  EXPECT_TRUE(::btool::util::fs::RemoveAll(dir, &err)) << "error: " << err;
}
