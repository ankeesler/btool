#include "runner.h"

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "node/node.h"

using ::testing::Ref;

class MockCallback : public ::btool::app::runner::Runner::Callback {
 public:
  MOCK_METHOD1(OnRun, void(const ::btool::node::Node &));
};

TEST(Runner, Success) {
  MockCallback mcb;
  ::btool::node::Node n("echo");
  EXPECT_CALL(mcb, OnRun(Ref(n)));

  ::btool::app::runner::Runner r(&mcb);
  EXPECT_TRUE(r.Run(n));
}

TEST(Runner, Failure) {
  // Note: this mock will be seen as leaked because when we fail to run a
  // child process ("this-binary-does-not-exist") then the child will exit(3)
  // and thus the mock will not be freed in that process. It isn't leaked in
  // the parent process though!
  MockCallback mcb;
  ::btool::node::Node n("this-binary-does-not-exist");
  EXPECT_CALL(mcb, OnRun(Ref(n)));

  ::btool::app::runner::Runner r(&mcb);
  EXPECT_FALSE(r.Run(n));
}
