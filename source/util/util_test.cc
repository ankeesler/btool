#include "util/util.h"

#include "gtest/gtest.h"

TEST(Hex, Success) {
  EXPECT_EQ("74756e61464953484d61526c496e",
            ::btool::util::Hex("tunaFISHMaRlIn"));
}
