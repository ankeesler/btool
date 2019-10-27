#include "util/string/string.h"

#include "gtest/gtest.h"

TEST(String, HasSuffix) {
  EXPECT_TRUE(::btool::util::string::HasSuffix("abc123", "123"));
  EXPECT_FALSE(::btool::util::string::HasSuffix("abc123", "abc"));
  EXPECT_TRUE(::btool::util::string::HasSuffix("abc123", "abc123"));
}
