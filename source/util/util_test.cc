#include "util/util.h"

#include "gtest/gtest.h"

TEST(Hex, Success) {
  unsigned char bytes[] = {
      0xAB,
      0x12,
      0xDE,
      0x34,
  };
  EXPECT_EQ("ab12de34",
            ::btool::util::Hex(bytes, sizeof(bytes) / sizeof(bytes[0])));
}
