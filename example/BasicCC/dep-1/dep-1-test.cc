#include "dep-1.h"

#include "gtest/gtest.h"

#include "dep-0/dep-0.h"

TEST(Hey, Sup) {
  EXPECT_EQ(1, 1);
  EXPECT_EQ(1, 1);
  EXPECT_NE(0, 1);
}
