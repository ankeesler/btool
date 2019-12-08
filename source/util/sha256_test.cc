#include "util/sha256.h"

#include <sstream>

#include "gtest/gtest.h"

TEST(SHA256, Yeah) {
  std::stringstream ss{"hey what's up"};
  EXPECT_EQ("19d4a0529e23484180d04cb47c6b914912a84a0c6521798c06f68e6363dbc65b",
            ::btool::util::SHA256(&ss));
}
