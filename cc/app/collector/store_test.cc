#include "store.h"

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "node/node.h"

using ::testing::InSequence;

TEST(Store, Basic) {
  ::btool::app::collector::Store s;
  ::btool::node::Node *b = s.Put("b");
  ::btool::node::Node *a = s.Put("a");
  a->dependencies()->push_back(b);

  EXPECT_EQ(a, s.Get("a"));
  EXPECT_EQ(b, s.Get("a")->dependencies()->at(0));
  EXPECT_EQ(b, s.Get("b"));
  EXPECT_EQ(nullptr, s.Get("c"));

  EXPECT_EQ(b, s.Put("b"));
}
