#include "util/flags.h"

#include <map>
#include <sstream>

#include "gtest/gtest.h"

TEST(Flags, Yeah) {
  ::btool::util::Flags f;

  bool a = false, b = true, c = false;
  f.Bool("a", "some-a-description", &a);
  f.Bool("b", "some-b-description", &b);
  f.Bool("c", "some-c-description", &c);
  std::string as = "a", bs = "b";
  f.String("as", "some-as-description", &as);
  f.String("bs", "some-bs-description", &bs);

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

  std::stringstream ss;
  f.Usage(&ss);
  const std::string ex =
      "  -a\n"
      "    some-a-description\n"
      "  -b\n"
      "    some-b-description\n"
      "  -c\n"
      "    some-c-description\n"
      "  -as\n"
      "    some-as-description\n"
      "  -bs\n"
      "    some-bs-description\n";
  EXPECT_EQ(ex, ss.str());
}
