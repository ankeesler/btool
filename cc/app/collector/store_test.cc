#include "store.h"

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "node/node.h"

using ::testing::InSequence;

class MockListener : public ::btool::app::collector::Store::Listener {
 public:
  MOCK_METHOD2(OnSet,
               void(::btool::app::collector::Store *, const std::string &));
};

TEST(Store, Basic) {
  ::btool::app::collector::Store s;
  ::btool::node::Node *b = s.Put("b");
  ::btool::node::Node *a = s.Put("a");
  a->AddDep(b);

  EXPECT_EQ(a, s.Get("a"));
  EXPECT_EQ(b, s.Get("a")->Deps()[0]);
  EXPECT_EQ(b, s.Get("b"));
  EXPECT_EQ(nullptr, s.Get("c"));

  EXPECT_EQ(b, s.Put("b"));
}

TEST(Store, Listener) {
  ::btool::app::collector::Store store;

  MockListener l0;
  MockListener l1;

  InSequence s;
  EXPECT_CALL(l0, OnSet(&store, "a"));
  EXPECT_CALL(l1, OnSet(&store, "a"));

  EXPECT_CALL(l0, OnSet(&store, "b"));
  EXPECT_CALL(l1, OnSet(&store, "b"));

  EXPECT_CALL(l0, OnSet(&store, "b"));
  EXPECT_CALL(l1, OnSet(&store, "b"));

  EXPECT_CALL(l0, OnSet(&store, "a"));
  EXPECT_CALL(l1, OnSet(&store, "a"));

  store.Listen(&l0);
  store.Listen(&l1);

  ::btool::node::Node *a = store.Put("a");
  ::btool::node::Node *b = store.Put("b");

  store.Set(b);
  store.Set(a);
}
