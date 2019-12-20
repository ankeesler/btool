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
  MOCK_METHOD1(RemoveAll, void(const std::string &));
};

class CleanerTest : public ::btool::node::testing::NodeTest {};

TEST_F(CleanerTest, Success) {
  InSequence s;

  MockRemoveAller mra;
  EXPECT_CALL(mra, RemoveAll("d"));
  EXPECT_CALL(mra, RemoveAll("c"));
  EXPECT_CALL(mra, RemoveAll("b"));
  EXPECT_CALL(mra, RemoveAll("a"));

  ::btool::app::cleaner::Cleaner cleaner(&mra);

  cleaner.Clean(a_);
}
