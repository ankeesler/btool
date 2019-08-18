#include "dep-0.h"

//#include "gtest/gtest.h"

using testing::Eq;

TEST(ProgramTest, Equals) {
  EXPECT_THAT(0, Eq(0));
}
