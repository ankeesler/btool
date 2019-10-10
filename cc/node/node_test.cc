#include "node.h"

#include <sstream>
#include <vector>

// workaround for bug-00
#include "gmock/gmock.h"
#include "gtest/gtest.h"

// workaround for bug-00
#include "core/err.h"
#include "node/testing/node.h"

class NodeTest : public ::btool::node::testing::NodeTest {};

TEST(Node, First) {
  ::btool::node::Node a("a");
  EXPECT_EQ("a", a.Name());
}

TEST_F(NodeTest, Print) {
  std::string ex = "a\n. b\n. . c\n. . . d\n. c\n. . d\n";
  std::stringstream ss;
  a_.String(&ss);
  EXPECT_EQ(ex, ss.str());
}

TEST_F(NodeTest, Visit) {
  std::vector<const ::btool::node::Node *> ex{&d_, &c_, &b_, &a_};
  std::vector<const ::btool::node::Node *> visited;
  a_.Visit(
      [&visited](const ::btool::node::Node *vn) { visited.push_back(vn); });
  EXPECT_EQ(ex, visited);
}

TEST_F(NodeTest, Deps) {
  auto deps = a_.Deps();
  EXPECT_EQ(2UL, deps.size());
  EXPECT_EQ("b", deps[0]->Name());
  EXPECT_EQ("c", deps[1]->Name());
}
