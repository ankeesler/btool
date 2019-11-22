#include "cleaner.h"

#include <string>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

// workaround for bug-00
#include "err.h"
#include "node/node.h"
#include "node/testing/node.h"

using ::testing::_;
using ::testing::InSequence;
using ::testing::Return;

class MockRemoveAller : public ::btool::app::cleaner::Cleaner::RemoveAller {
 public:
  MOCK_METHOD2(RemoveAll, bool(const std::string &, std::string *));
};

class CleanerTest : public ::btool::node::testing::NodeTest {};

TEST_F(CleanerTest, Success) {
  InSequence s;

  MockRemoveAller mra;
  EXPECT_CALL(mra, RemoveAll("d", _)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("c", _)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("b", _)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("a", _)).WillOnce(Return(true));

  ::btool::app::cleaner::Cleaner cleaner(&mra);

  std::string err;
  EXPECT_TRUE(cleaner.Clean(a_, &err)) << "error: " << err;
}

TEST_F(CleanerTest, Failure) {
  InSequence s;

  MockRemoveAller mra;
  EXPECT_CALL(mra, RemoveAll("d", _)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("c", _)).WillOnce(Return(false));

  ::btool::app::cleaner::Cleaner cleaner(&mra);

  std::string err;
  EXPECT_FALSE(cleaner.Clean(a_, &err));
}
