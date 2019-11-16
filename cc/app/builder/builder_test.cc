#include "builder.h"

#include <string>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "err.h"
#include "node/node.h"
#include "node/testing/node.h"

using ::testing::_;
using ::testing::InSequence;
using ::testing::Ref;
using ::testing::Return;

class MockCurrenter : public ::btool::app::builder::Builder::Currenter {
 public:
  MOCK_METHOD1(Current, ::btool::Err<bool>(const ::btool::node::Node &));
};

class BuilderTest : public ::btool::node::testing::NodeTest {};

TEST_F(BuilderTest, BuildAll) {
  InSequence s;
  MockCurrenter mc;
  EXPECT_CALL(mc, Current(Ref(d_)))
      .WillOnce(Return(::btool::Err<bool>::Success(false)));
  EXPECT_CALL(dr_, Resolve(Ref(d_)))
      .WillOnce(Return(::btool::VoidErr::Success()));

  EXPECT_CALL(mc, Current(Ref(c_)))
      .WillOnce(Return(::btool::Err<bool>::Success(false)));
  EXPECT_CALL(cr_, Resolve(Ref(c_)))
      .WillOnce(Return(::btool::VoidErr::Success()));

  EXPECT_CALL(mc, Current(Ref(b_)))
      .WillOnce(Return(::btool::Err<bool>::Success(false)));
  EXPECT_CALL(br_, Resolve(Ref(b_)))
      .WillOnce(Return(::btool::VoidErr::Success()));

  EXPECT_CALL(mc, Current(Ref(a_)))
      .WillOnce(Return(::btool::Err<bool>::Success(false)));
  EXPECT_CALL(ar_, Resolve(Ref(a_)))
      .WillOnce(Return(::btool::VoidErr::Success()));

  ::btool::app::builder::Builder b(&mc);
  EXPECT_EQ(::btool::VoidErr::Success(), b.Build(a_));
}

TEST_F(BuilderTest, UpToDate) {
  InSequence s;
  MockCurrenter mc;
  EXPECT_CALL(mc, Current(Ref(d_)))
      .WillOnce(Return(::btool::Err<bool>::Success(true)));

  EXPECT_CALL(mc, Current(Ref(c_)))
      .WillOnce(Return(::btool::Err<bool>::Success(true)));

  EXPECT_CALL(mc, Current(Ref(b_)))
      .WillOnce(Return(::btool::Err<bool>::Success(true)));

  EXPECT_CALL(mc, Current(Ref(a_)))
      .WillOnce(Return(::btool::Err<bool>::Success(true)));

  ::btool::app::builder::Builder b(&mc);
  EXPECT_EQ(::btool::VoidErr::Success(), b.Build(a_));
}

TEST_F(BuilderTest, Some) {
  InSequence s;
  MockCurrenter mc;
  EXPECT_CALL(mc, Current(Ref(d_)))
      .WillOnce(Return(::btool::Err<bool>::Success(true)));

  EXPECT_CALL(mc, Current(Ref(c_)))
      .WillOnce(Return(::btool::Err<bool>::Success(false)));
  EXPECT_CALL(cr_, Resolve(Ref(c_)))
      .WillOnce(Return(::btool::VoidErr::Success()));

  EXPECT_CALL(mc, Current(Ref(b_)))
      .WillOnce(Return(::btool::Err<bool>::Success(false)));
  EXPECT_CALL(br_, Resolve(Ref(b_)))
      .WillOnce(Return(::btool::VoidErr::Success()));

  EXPECT_CALL(mc, Current(Ref(a_)))
      .WillOnce(Return(::btool::Err<bool>::Success(true)));

  ::btool::app::builder::Builder b(&mc);
  EXPECT_EQ(::btool::VoidErr::Success(), b.Build(a_));
}
