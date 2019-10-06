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

  std::string ex = "0\n. 1\n. . 2\n. . . 3\n. 2\n. . 3\n";
  std::stringstream ss;
  nodes->Get("0")->String(&ss);
  EXPECT_EQ(ex, ss.str());
}

TEST(Node, Visit) {
  auto nodes = ::btool::node::testing::Nodes0123();

  std::vector<const ::btool::node::Node *> ex{nodes->Get("3"), nodes->Get("2"),
                                              nodes->Get("1"), nodes->Get("0")};
  std::vector<const ::btool::node::Node *> visited;
  nodes->Get("0")->Visit(
      [&visited](const ::btool::node::Node *vn) { visited.push_back(vn); });
  EXPECT_EQ(ex, visited);
}

TEST(Node, Deps) {
  auto nodes = ::btool::node::testing::Nodes0123();

  auto deps = nodes->Get("0")->Deps();
  EXPECT_EQ(2, deps.size());
  EXPECT_EQ("1", deps[0]->Name());
  EXPECT_EQ("2", deps[1]->Name());
}
