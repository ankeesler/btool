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

TEST(ErrorTest, Equals) {
  btool::Error e0 = ::btool::Error::Create("some message");
  btool::Error e1 = ::btool::Error::Create("some other message");
  btool::Error e2 = ::btool::Error::Create("some message");
  btool::Error e3 = ::btool::Error::Success();

  EXPECT_EQ(e0, e0);
  EXPECT_NE(e0, e1);
  EXPECT_EQ(e0, e2);
  EXPECT_NE(e0, e3);

  EXPECT_EQ(e1, e1);
  EXPECT_NE(e1, e2);
  EXPECT_NE(e1, e3);

  EXPECT_EQ(e2, e2);
  EXPECT_NE(e2, e3);

  EXPECT_EQ(e3, e3);
}
