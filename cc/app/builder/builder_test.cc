#include "builder.h"

#include <string>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "err.h"
#include "node/node.h"
#include "node/testing/node.h"

using ::testing::_;
using ::testing::DoAll;
using ::testing::InSequence;
using ::testing::Ref;
using ::testing::Return;
using ::testing::SetArgPointee;

class MockCurrenter : public ::btool::app::builder::Builder::Currenter {
 public:
  MOCK_METHOD3(Current,
               bool(const ::btool::node::Node &, bool *, std::string *));
};

class BuilderTest : public ::btool::node::testing::NodeTest {};

TEST_F(BuilderTest, BuildAll) {
  InSequence s;
  MockCurrenter mc;
  EXPECT_CALL(mc, Current(Ref(d_), _, _))
      .WillOnce(DoAll(SetArgPointee<1>(false), Return(true)));
  EXPECT_CALL(dr_, Resolve(Ref(d_)))
      .WillOnce(Return(::btool::VoidErr::Success()));

  EXPECT_CALL(mc, Current(Ref(c_), _, _))
      .WillOnce(DoAll(SetArgPointee<1>(false), Return(true)));
  EXPECT_CALL(cr_, Resolve(Ref(c_)))
      .WillOnce(Return(::btool::VoidErr::Success()));

  EXPECT_CALL(mc, Current(Ref(b_), _, _))
      .WillOnce(DoAll(SetArgPointee<1>(false), Return(true)));
  EXPECT_CALL(br_, Resolve(Ref(b_)))
      .WillOnce(Return(::btool::VoidErr::Success()));

  EXPECT_CALL(mc, Current(Ref(a_), _, _))
      .WillOnce(DoAll(SetArgPointee<1>(false), Return(true)));
  EXPECT_CALL(ar_, Resolve(Ref(a_)))
      .WillOnce(Return(::btool::VoidErr::Success()));

  ::btool::app::builder::Builder b(&mc);
  std::string ret_err;
  EXPECT_TRUE(b.Build(a_, &ret_err)) << "error: " << ret_err;
}

TEST_F(BuilderTest, UpToDate) {
  InSequence s;
  MockCurrenter mc;
  EXPECT_CALL(mc, Current(Ref(d_), _, _))
      .WillOnce(DoAll(SetArgPointee<1>(true), Return(true)));

  EXPECT_CALL(mc, Current(Ref(c_), _, _))
      .WillOnce(DoAll(SetArgPointee<1>(true), Return(true)));

  EXPECT_CALL(mc, Current(Ref(b_), _, _))
      .WillOnce(DoAll(SetArgPointee<1>(true), Return(true)));

  EXPECT_CALL(mc, Current(Ref(a_), _, _))
      .WillOnce(DoAll(SetArgPointee<1>(true), Return(true)));

  ::btool::app::builder::Builder b(&mc);
  std::string ret_err;
  EXPECT_TRUE(b.Build(a_, &ret_err)) << "error: " << ret_err;
}

TEST_F(BuilderTest, Some) {
  InSequence s;
  MockCurrenter mc;
  EXPECT_CALL(mc, Current(Ref(d_), _, _))
      .WillOnce(DoAll(SetArgPointee<1>(true), Return(true)));

  EXPECT_CALL(mc, Current(Ref(c_), _, _))
      .WillOnce(DoAll(SetArgPointee<1>(false), Return(true)));
  EXPECT_CALL(cr_, Resolve(Ref(c_)))
      .WillOnce(Return(::btool::VoidErr::Success()));

  EXPECT_CALL(mc, Current(Ref(b_), _, _))
      .WillOnce(DoAll(SetArgPointee<1>(false), Return(true)));
  EXPECT_CALL(br_, Resolve(Ref(b_)))
      .WillOnce(Return(::btool::VoidErr::Success()));

  EXPECT_CALL(mc, Current(Ref(a_), _, _))
      .WillOnce(DoAll(SetArgPointee<1>(true), Return(true)));

  ::btool::app::builder::Builder b(&mc);
  std::string ret_err;
  EXPECT_TRUE(b.Build(a_, &ret_err)) << "error: " << ret_err;
}
