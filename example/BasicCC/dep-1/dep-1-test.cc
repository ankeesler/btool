#include "dep-1.h"

#include "gtest/gtest.h"

#include "dep-0/dep-0.h"

TEST(Hey, Sup) {
  EXPECT_EQ(1, 1);
}

int main(int argc, char **argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
