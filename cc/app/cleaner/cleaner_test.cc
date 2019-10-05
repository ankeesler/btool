#include "cleaner.h"

#include <string>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "node/node.h"

using ::testing::InSequence;
using ::testing::Return;

class MockRemoveAller : public ::btool::app::cleaner::Cleaner::RemoveAller {
 public:
  MOCK_METHOD2(RemoveAll, bool(const std::string &, std::string *));
};

TEST(Cleaner, Success) {
  InSequence s;

  std::string err;
  MockRemoveAller mra;
  EXPECT_CALL(mra, RemoveAll("d", &err)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("c", &err)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("b", &err)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("a", &err)).WillOnce(Return(true));

  ::btool::app::cleaner::Cleaner cleaner(&mra);

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

  EXPECT_TRUE(cleaner.Clean(a, &err)) << err;
}

TEST(Cleaner, Failure) {
  InSequence s;

  std::string err;
  MockRemoveAller mra;
  EXPECT_CALL(mra, RemoveAll("d", &err)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("c", &err)).WillOnce(Return(false));

  ::btool::app::cleaner::Cleaner cleaner(&mra);

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

  EXPECT_FALSE(cleaner.Clean(a, &err)) << err;
}
