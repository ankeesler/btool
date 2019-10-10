#include "cleaner.h"

#include <string>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

// workaround for bug-00
#include "core/err.h"
#include "node/node.h"
#include "node/testing/node.h"

using ::testing::InSequence;
using ::testing::Return;

class MockRemoveAller : public ::btool::app::cleaner::Cleaner::RemoveAller {
 public:
  MOCK_METHOD2(RemoveAll, bool(const std::string &, std::string *));
};

class CleanerTest : public ::btool::node::testing::NodeTest {};

TEST_F(CleanerTest, Success) {
  InSequence s;

  std::string err;
  MockRemoveAller mra;
  EXPECT_CALL(mra, RemoveAll("d", &err)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("c", &err)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("b", &err)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("a", &err)).WillOnce(Return(true));

  ::btool::app::cleaner::Cleaner cleaner(&mra);

  EXPECT_TRUE(cleaner.Clean(a_, &err)) << err;
}

TEST_F(CleanerTest, Failure) {
  InSequence s;

  std::string err;
  MockRemoveAller mra;
  EXPECT_CALL(mra, RemoveAll("d", &err)).WillOnce(Return(true));
  EXPECT_CALL(mra, RemoveAll("c", &err)).WillOnce(Return(false));

  ::btool::app::cleaner::Cleaner cleaner(&mra);

  EXPECT_FALSE(cleaner.Clean(a_, &err)) << err;
}
