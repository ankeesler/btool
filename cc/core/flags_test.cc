#include "flags.h"

#include <map>

#include "gtest/gtest.h"

TEST(Flags, Bool) {
  ::btool::core::Flags f;

  bool a = false, b = true, c = false;
  f.Bool("a", &a);
  f.Bool("b", &b);
  f.Bool("c", &c);

  std::string err;
  const char *argv[] = {
      "-a",
      "-c",
  };
  int argc = sizeof(argv) / sizeof(argv[0]);
  bool success = f.Parse(argc, argv, &err);
  EXPECT_TRUE(success);
  EXPECT_TRUE(a);
  EXPECT_FALSE(b);
  EXPECT_TRUE(c);
}