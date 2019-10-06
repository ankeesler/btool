#include "node/store.h"

#include "gtest/gtest.h"

#include "node/node.h"

TEST(Store, Basic) {
  ::btool::node::Store s;
  ::btool::node::Node *b = s.Create("b");
  ::btool::node::Node *a = s.Create("a");
  a->AddDep(b);

  EXPECT_EQ(a, s.Get("a"));
  EXPECT_EQ(b, s.Get("a")->Deps()[0]);
  EXPECT_EQ(b, s.Get("b"));
  EXPECT_EQ(nullptr, s.Get("c"));
}
