#include "node.h"

#include <sstream>

#include "gtest/gtest.h"

TEST(Node, First) {
  ::btool::node::Node a("a");
  EXPECT_EQ("a", a.Name());
}

TEST(Node, Print) {
  // a -> b, c
  // b -> c
  // c -> d
  // d
  ::btool::node::Node d("d");
  ::btool::node::Node c("c");
  c.AddDep(&d);
  ::btool::node::Node b("b");
  b.AddDep(&c);
  ::btool::node::Node a("a");
  a.AddDep(&b);
  a.AddDep(&c);

  std::string ex = "a\n. b\n. . c\n. . . d\n. c\n. . d\n";
  std::stringstream ss;
  a.String(&ss);
  EXPECT_EQ(ex, ss.str());
}

int main(int argc, char *argv[]) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
