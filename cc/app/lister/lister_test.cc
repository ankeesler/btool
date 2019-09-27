#include "lister.h"

#include <sstream>

#include "gtest/gtest.h"

#include "node/node.h"

TEST(Lister, List) {
  std::stringstream ss;
  ::btool::app::lister::Lister l(&ss);

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

  l.List(a);

  std::string ex = "a\n. b\n. . c\n. . . d\n. c\n. . d\n";
  EXPECT_EQ(ex, ss.str());
}
