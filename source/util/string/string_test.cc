#include "util/string/string.h"

#include "gtest/gtest.h"

TEST(String, HasPrefix) {
  EXPECT_FALSE(::btool::util::string::HasPrefix("abc123", "123"));
  EXPECT_FALSE(::btool::util::string::HasPrefix("abc123", "123abc456def"));
  EXPECT_TRUE(::btool::util::string::HasPrefix("abc123", "abc"));
  EXPECT_TRUE(::btool::util::string::HasPrefix("abc123", "abc123"));
}

TEST(String, HasSuffix) {
  EXPECT_TRUE(::btool::util::string::HasSuffix("abc123", "123"));
  EXPECT_FALSE(::btool::util::string::HasSuffix("abc123", "abc"));
  EXPECT_FALSE(::btool::util::string::HasSuffix("abc123", "abc123def456"));
  EXPECT_TRUE(::btool::util::string::HasSuffix("abc123", "abc123"));
}

TEST(String, Replace) {
  EXPECT_EQ("foo.o", ::btool::util::string::Replace("foo.cc", ".cc", ".o"));
  EXPECT_EQ("foo.cc", ::btool::util::string::Replace("foo.o", ".o", ".cc"));
  EXPECT_EQ("foo.b", ::btool::util::string::Replace("foo.a", "a", "b"));
}
