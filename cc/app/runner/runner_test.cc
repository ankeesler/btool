#include "runner.h"

#include "gtest/gtest.h"

#include "node/node.h"

TEST(Runner, Run) {
  ::btool::app::runner::Runner r;
  ::btool::node::Node n("a");
  EXPECT_TRUE(r.Run(n, nullptr));
}
