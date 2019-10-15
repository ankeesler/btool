#include "flags.h"

#include <map>

#include "gtest/gtest.h"

TEST(Flags, Bool) {
  ::btool::util::Flags f;

  bool a = false, b = true, c = false;
  f.Bool("a", &a);
  f.Bool("b", &b);
  f.Bool("c", &c);
  std::string as = "a", bs = "b";
  f.String("as", &as);
  f.String("bs", &bs);

  std::string err;
  const char *argv[] = {
      "-a",
      "-c",
      "-as",
      "tuna",
  };
  int argc = sizeof(argv) / sizeof(argv[0]);
  bool success = f.Parse(argc, argv, &err);
  EXPECT_TRUE(success);
  EXPECT_TRUE(a);
  EXPECT_FALSE(b);
  EXPECT_TRUE(c);
  EXPECT_EQ("tuna", as);
  EXPECT_EQ("b", bs);
}
