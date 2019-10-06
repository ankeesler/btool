#include "cleaner.h"

#include <string>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "node/node.h"
#include "node/testing/node.h"

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

  auto nodes = ::btool::node::testing::Nodes0123();

  EXPECT_TRUE(cleaner.Clean(*nodes->at(0), &err)) << err;

  ::btool::node::testing::Free(nodes);
}

TEST(Cleaner, Failure) {
  InSequence s;

  std::string err;
  MockRemoveAller mra;
  EXPECT_CALL(mra, RemoveAll("d", &err)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("c", &err)).WillOnce(Return(false));

  ::btool::app::cleaner::Cleaner cleaner(&mra);

  auto nodes = ::btool::node::testing::Nodes0123();

  EXPECT_FALSE(cleaner.Clean(*nodes->at(0), &err)) << err;

  ::btool::node::testing::Free(nodes);
}
