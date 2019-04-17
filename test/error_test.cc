#include "gtest/gtest.h"

#include "error.h"

TEST(ErrorTest, Create) {
  btool::Error e = ::btool::Error::Create("some message");
  EXPECT_TRUE(e.Exists());
  EXPECT_EQ(e.Message(), "some message");
}

TEST(ErrorTest, Success) {
  btool::Error e = ::btool::Error::Success();
  EXPECT_FALSE(e.Exists());
  e.Message(); // nothing bad should happen
}
