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
  MOCK_METHOD1(RemoveAll, ::btool::core::VoidErr(const std::string &));
};

class CleanerTest : public ::btool::node::testing::NodeTest {};

TEST_F(CleanerTest, Success) {
  InSequence s;

  MockRemoveAller mra;
  EXPECT_CALL(mra, RemoveAll("d"))
      .WillOnce(Return(::btool::core::VoidErr::Success()));
  EXPECT_CALL(mra, RemoveAll("c"))
      .WillOnce(Return(::btool::core::VoidErr::Success()));
  EXPECT_CALL(mra, RemoveAll("b"))
      .WillOnce(Return(::btool::core::VoidErr::Success()));
  EXPECT_CALL(mra, RemoveAll("a"))
      .WillOnce(Return(::btool::core::VoidErr::Success()));

  ::btool::app::cleaner::Cleaner cleaner(&mra);

  EXPECT_FALSE(cleaner.Clean(a_));
}

TEST_F(CleanerTest, Failure) {
  InSequence s;

  MockRemoveAller mra;
  EXPECT_CALL(mra, RemoveAll("d"))
      .WillOnce(Return(::btool::core::VoidErr::Success()));
  EXPECT_CALL(mra, RemoveAll("c"))
      .WillOnce(Return(::btool::core::VoidErr::Failure("eh")));

  ::btool::app::cleaner::Cleaner cleaner(&mra);

  EXPECT_TRUE(cleaner.Clean(a_));
}
