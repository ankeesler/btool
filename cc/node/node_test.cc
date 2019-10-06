#include "node.h"

#include <sstream>
#include <vector>

#include "gtest/gtest.h"

#include "node/testing/node.h"

TEST(Node, First) {
  ::btool::node::Node a("a");
  EXPECT_EQ("a", a.Name());
}

TEST(Node, Print) {
  auto nodes = ::btool::node::testing::Nodes0123();

  std::string ex = "a\n. b\n. . c\n. . . d\n. c\n. . d\n";
  std::stringstream ss;
  nodes->at(0)->String(&ss);
  EXPECT_EQ(ex, ss.str());

  ::btool::node::testing::Free(nodes);
}

TEST(Node, Visit) {
  auto nodes = ::btool::node::testing::Nodes0123();

  std::vector<const ::btool::node::Node *> ex{nodes->at(3), nodes->at(2),
                                              nodes->at(1), nodes->at(0)};
  std::vector<const ::btool::node::Node *> visited;
  nodes->at(0)->Visit(
      [&visited](const ::btool::node::Node *vn) { visited.push_back(vn); });
  EXPECT_EQ(ex, visited);

  ::btool::node::testing::Free(nodes);
}

TEST(Node, Deps) {
  auto nodes = ::btool::node::testing::Nodes0123();

  auto deps = nodes->at(0)->Deps();
  EXPECT_EQ(2, deps.size());
  EXPECT_EQ("b", deps[0]->Name());
  EXPECT_EQ("c", deps[1]->Name());

  ::btool::node::testing::Free(nodes);
}
