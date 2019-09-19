#include "node.h"

#include "gtest/gtest.h"

TEST(Node, First) {
  ::btool::node::Node a("a");
  EXPECT_EQ("a", a.Name());
}

int main(int argc, char *argv[]) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
