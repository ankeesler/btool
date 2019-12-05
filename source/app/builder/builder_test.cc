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
using ::testing::StrictMock;

class MockCurrenter : public ::btool::app::builder::Builder::Currenter {
 public:
  MOCK_METHOD1(Current, bool(const ::btool::node::Node &));
};

class MockCallback : public ::btool::app::builder::Builder::Callback {
 public:
  MOCK_METHOD2(OnResolve, void(const ::btool::node::Node &, bool));
};

class BuilderTest : public ::btool::node::testing::NodeTest {};

TEST_F(BuilderTest, BuildAll) {
  InSequence s;
  StrictMock<MockCurrenter> mcu;
  StrictMock<MockCallback> mca;
  EXPECT_CALL(mcu, Current(Ref(d_))).WillOnce(Return(false));
  EXPECT_CALL(mca, OnResolve(Ref(d_), false));
  EXPECT_CALL(dr_, Resolve(Ref(d_)));

  EXPECT_CALL(mcu, Current(Ref(c_))).WillOnce(Return(false));
  EXPECT_CALL(mca, OnResolve(Ref(c_), false));
  EXPECT_CALL(cr_, Resolve(Ref(c_)));

  EXPECT_CALL(mcu, Current(Ref(b_))).WillOnce(Return(false));
  EXPECT_CALL(mca, OnResolve(Ref(b_), false));
  EXPECT_CALL(br_, Resolve(Ref(b_)));

  EXPECT_CALL(mcu, Current(Ref(a_))).WillOnce(Return(false));
  EXPECT_CALL(mca, OnResolve(Ref(a_), false));
  EXPECT_CALL(ar_, Resolve(Ref(a_)));

  ::btool::app::builder::Builder b(&mcu, &mca);
  b.Build(a_);
}

TEST_F(BuilderTest, UpToDate) {
  InSequence s;
  StrictMock<MockCurrenter> mcu;
  StrictMock<MockCallback> mca;
  EXPECT_CALL(mcu, Current(Ref(d_))).WillOnce(Return(true));
  EXPECT_CALL(mca, OnResolve(Ref(d_), true));

  EXPECT_CALL(mcu, Current(Ref(c_))).WillOnce(Return(true));
  EXPECT_CALL(mca, OnResolve(Ref(c_), true));

  EXPECT_CALL(mcu, Current(Ref(b_))).WillOnce(Return(true));
  EXPECT_CALL(mca, OnResolve(Ref(b_), true));

  EXPECT_CALL(mcu, Current(Ref(a_))).WillOnce(Return(true));
  EXPECT_CALL(mca, OnResolve(Ref(a_), true));

  ::btool::app::builder::Builder b(&mcu, &mca);
  b.Build(a_);
}

TEST_F(BuilderTest, Some) {
  InSequence s;
  StrictMock<MockCurrenter> mcu;
  StrictMock<MockCallback> mca;
  EXPECT_CALL(mcu, Current(Ref(d_))).WillOnce(Return(true));
  EXPECT_CALL(mca, OnResolve(Ref(d_), true));

  EXPECT_CALL(mcu, Current(Ref(c_))).WillOnce(Return(false));
  EXPECT_CALL(mca, OnResolve(Ref(c_), false));
  EXPECT_CALL(cr_, Resolve(Ref(c_)));

  EXPECT_CALL(mcu, Current(Ref(b_))).WillOnce(Return(false));
  EXPECT_CALL(mca, OnResolve(Ref(b_), false));
  EXPECT_CALL(br_, Resolve(Ref(b_)));

  EXPECT_CALL(mcu, Current(Ref(a_))).WillOnce(Return(true));
  EXPECT_CALL(mca, OnResolve(Ref(a_), true));

  ::btool::app::builder::Builder b(&mcu, &mca);
  b.Build(a_);
}
