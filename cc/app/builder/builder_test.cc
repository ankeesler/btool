#include "builder.h"

#include <string>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "core/err.h"
#include "node/node.h"
#include "node/testing/node.h"

using ::testing::_;
using ::testing::InSequence;
using ::testing::Ref;
using ::testing::Return;

class MockCurrenter : public ::btool::app::builder::Builder::Currenter {
 public:
  MOCK_METHOD1(Current, ::btool::core::Err<bool>(const ::btool::node::Node &));
};

TEST(Builder, BuildAll) {
  auto nodes = ::btool::node::testing::Nodes0123();

  ::btool::node::testing::MockResolver mr0;
  nodes->Get("0")->SetResolver(&mr0);
  ::btool::node::testing::MockResolver mr1;
  nodes->Get("1")->SetResolver(&mr1);
  ::btool::node::testing::MockResolver mr2;
  nodes->Get("2")->SetResolver(&mr2);
  ::btool::node::testing::MockResolver mr3;
  nodes->Get("3")->SetResolver(&mr3);

  InSequence s;
  MockCurrenter mc;
  std::string err;
  EXPECT_CALL(mc, Current(Ref(*nodes->Get("3"))))
      .WillOnce(Return(::btool::core::Err<bool>::Success(false)));
  EXPECT_CALL(mr3, Resolve(Ref(*nodes->Get("3"))))
      .WillOnce(Return(::btool::core::VoidErr::Success()));

  EXPECT_CALL(mc, Current(Ref(*nodes->Get("2"))))
      .WillOnce(Return(::btool::core::Err<bool>::Success(false)));
  EXPECT_CALL(mr2, Resolve(Ref(*nodes->Get("2"))))
      .WillOnce(Return(::btool::core::VoidErr::Success()));

  EXPECT_CALL(mc, Current(Ref(*nodes->Get("1"))))
      .WillOnce(Return(::btool::core::Err<bool>::Success(false)));
  EXPECT_CALL(mr1, Resolve(Ref(*nodes->Get("1"))))
      .WillOnce(Return(::btool::core::VoidErr::Success()));

  EXPECT_CALL(mc, Current(Ref(*nodes->Get("0"))))
      .WillOnce(Return(::btool::core::Err<bool>::Success(false)));
  EXPECT_CALL(mr0, Resolve(Ref(*nodes->Get("0"))))
      .WillOnce(Return(::btool::core::VoidErr::Success()));

  ::btool::app::builder::Builder b(&mc);
  EXPECT_EQ(::btool::core::VoidErr::Success(), b.Build(*nodes->Get("0")));
}
