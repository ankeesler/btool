#include "cmd.h"

#include <sstream>

#include "gtest/gtest.h"

TEST(Cmd, Pass) {
  ::btool::core::Cmd c("echo");

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
  ::btool::core::Cmd c("cat");

  c.Arg("this/file/does/not/exist");

  std::stringstream out, err;
  c.Stdout(&out);
  c.Stderr(&err);

  EXPECT_EQ(1, c.Run());

  EXPECT_EQ("", out.str());
  EXPECT_EQ("cat: this/file/does/not/exist: No such file or directory\n",
            err.str());
}
