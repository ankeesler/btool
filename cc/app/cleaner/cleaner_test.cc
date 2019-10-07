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
  EXPECT_CALL(mra, RemoveAll("3", &err)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("2", &err)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("1", &err)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("0", &err)).WillOnce(Return(true));

  ::btool::app::cleaner::Cleaner cleaner(&mra);

  auto nodes = ::btool::node::testing::Nodes0123();

  EXPECT_TRUE(cleaner.Clean(*nodes->Get("0"), &err)) << err;
}

TEST(Cleaner, Failure) {
  InSequence s;

  std::string err;
  MockRemoveAller mra;
  EXPECT_CALL(mra, RemoveAll("3", &err)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("2", &err)).WillOnce(Return(false));

  ::btool::app::cleaner::Cleaner cleaner(&mra);

  auto nodes = ::btool::node::testing::Nodes0123();

  EXPECT_FALSE(cleaner.Clean(*nodes->Get("0"), &err)) << err;
}
