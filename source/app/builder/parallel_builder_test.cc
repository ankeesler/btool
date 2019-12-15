#include "app/builder/parallel_builder.h"

#include <functional>
#include <queue>

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

class ParallelBuilderTest : public ::btool::node::testing::NodeTest {};

class WorkPool : public ::btool::app::builder::ParallelBuilder::WorkPool {
 public:
  void Submit(std::function<const ::btool::node::Node *()> work) override {
    q_.push(work());
  }

  const ::btool::node::Node *Receive(::btool::Err *err) override {
    auto n = q_.front();
    q_.pop();
    return n;
  }

 private:
  std::queue<const ::btool::node::Node *> q_;
};

class MockCurrenter : public ::btool::app::builder::ParallelBuilder::Currenter {
 public:
  MOCK_METHOD1(Current, bool(const ::btool::node::Node &));
};

class MockCallback : public ::btool::app::builder::ParallelBuilder::Callback {
 public:
  MOCK_METHOD2(OnPreResolve, void(const ::btool::node::Node &, bool));
  MOCK_METHOD2(OnPostResolve, void(const ::btool::node::Node &, bool));
};

TEST_F(ParallelBuilderTest, BuildAll) {
  InSequence s;

  WorkPool wp;
  StrictMock<MockCurrenter> mcu;
  StrictMock<MockCallback> mca;
  EXPECT_CALL(mcu, Current(Ref(d_))).WillOnce(Return(false));
  EXPECT_CALL(mca, OnPreResolve(Ref(d_), false));
  EXPECT_CALL(dr_, Resolve(Ref(d_)));
  EXPECT_CALL(mca, OnPostResolve(Ref(d_), false));

  EXPECT_CALL(mcu, Current(Ref(c_))).WillOnce(Return(false));
  EXPECT_CALL(mca, OnPreResolve(Ref(c_), false));
  EXPECT_CALL(cr_, Resolve(Ref(c_)));
  EXPECT_CALL(mca, OnPostResolve(Ref(c_), false));

  EXPECT_CALL(mcu, Current(Ref(b_))).WillOnce(Return(false));
  EXPECT_CALL(mca, OnPreResolve(Ref(b_), false));
  EXPECT_CALL(br_, Resolve(Ref(b_)));
  EXPECT_CALL(mca, OnPostResolve(Ref(b_), false));

  EXPECT_CALL(mcu, Current(Ref(a_))).WillOnce(Return(false));
  EXPECT_CALL(mca, OnPreResolve(Ref(a_), false));
  EXPECT_CALL(ar_, Resolve(Ref(a_)));
  EXPECT_CALL(mca, OnPostResolve(Ref(a_), false));

  ::btool::app::builder::ParallelBuilder pb(&wp, &mcu, &mca);
  pb.Build(a_);
}

TEST_F(ParallelBuilderTest, UpToDate) {
  InSequence s;

  WorkPool wp;
  StrictMock<MockCurrenter> mcu;
  StrictMock<MockCallback> mca;
  EXPECT_CALL(mcu, Current(Ref(d_))).WillOnce(Return(true));
  EXPECT_CALL(mca, OnPreResolve(Ref(d_), true));
  EXPECT_CALL(mca, OnPostResolve(Ref(d_), true));

  EXPECT_CALL(mcu, Current(Ref(c_))).WillOnce(Return(true));
  EXPECT_CALL(mca, OnPreResolve(Ref(c_), true));
  EXPECT_CALL(mca, OnPostResolve(Ref(c_), true));

  EXPECT_CALL(mcu, Current(Ref(b_))).WillOnce(Return(true));
  EXPECT_CALL(mca, OnPreResolve(Ref(b_), true));
  EXPECT_CALL(mca, OnPostResolve(Ref(b_), true));

  EXPECT_CALL(mcu, Current(Ref(a_))).WillOnce(Return(true));
  EXPECT_CALL(mca, OnPreResolve(Ref(a_), true));
  EXPECT_CALL(mca, OnPostResolve(Ref(a_), true));

  ::btool::app::builder::ParallelBuilder pb(&wp, &mcu, &mca);
  pb.Build(a_);
}

TEST_F(ParallelBuilderTest, Some) {
  InSequence s;

  WorkPool wp;
  StrictMock<MockCurrenter> mcu;
  StrictMock<MockCallback> mca;
  EXPECT_CALL(mcu, Current(Ref(d_))).WillOnce(Return(true));
  EXPECT_CALL(mca, OnPreResolve(Ref(d_), true));
  EXPECT_CALL(mca, OnPostResolve(Ref(d_), true));

  EXPECT_CALL(mcu, Current(Ref(c_))).WillOnce(Return(false));
  EXPECT_CALL(mca, OnPreResolve(Ref(c_), false));
  EXPECT_CALL(cr_, Resolve(Ref(c_)));
  EXPECT_CALL(mca, OnPostResolve(Ref(c_), false));

  EXPECT_CALL(mcu, Current(Ref(b_))).WillOnce(Return(false));
  EXPECT_CALL(mca, OnPreResolve(Ref(b_), false));
  EXPECT_CALL(br_, Resolve(Ref(b_)));
  EXPECT_CALL(mca, OnPostResolve(Ref(b_), false));

  EXPECT_CALL(mcu, Current(Ref(a_))).WillOnce(Return(true));
  EXPECT_CALL(mca, OnPreResolve(Ref(a_), true));
  EXPECT_CALL(mca, OnPostResolve(Ref(a_), true));

  ::btool::app::builder::ParallelBuilder pb(&wp, &mcu, &mca);
  pb.Build(a_);
}
