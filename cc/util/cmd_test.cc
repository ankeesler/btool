#include "util/cmd.h"

#include <sstream>

#include "gtest/gtest.h"

TEST(Cmd, Pass) {
  ::btool::util::Cmd c("echo");

  c.Arg("-n");
  c.Arg("tuna");
  c.Arg("fish");
  c.Arg("marlin");

  std::stringstream out, err;
  c.Stdout(&out);
  c.Stderr(&err);

  EXPECT_EQ(0, c.Run());

  EXPECT_EQ("tuna fish marlin", out.str());
  EXPECT_EQ("", err.str());
}

TEST(Cmd, Fail) {
  ::btool::util::Cmd c("cat");

  c.Arg("this/file/does/not/exist");

  std::stringstream out, err;
  c.Stdout(&out);
  c.Stderr(&err);

  EXPECT_EQ(1, c.Run());

  EXPECT_EQ("", out.str());
  EXPECT_EQ("cat: this/file/does/not/exist: No such file or directory\n",
            err.str());
}

TEST(Cmd, DoesNotExist) {
  ::btool::util::Cmd c("this-binary-does-not-exist");
  EXPECT_EQ(255, c.Run());
}
