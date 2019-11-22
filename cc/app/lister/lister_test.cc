#include "lister.h"

#include <sstream>
#include <string>

// workaround for bug-00
#include "gmock/gmock.h"
#include "gtest/gtest.h"

// workaround for bug-00
#include "err.h"
#include "node/node.h"
#include "node/testing/node.h"

using ::testing::_;

class ListerTest : public ::btool::node::testing::NodeTest {};

TEST_F(ListerTest, List) {
  std::stringstream ss;
  ::btool::app::lister::Lister l(&ss);

  std::string err;
  EXPECT_TRUE(l.List(a_, &err)) << "error: " << err;

  std::string ex = "a\n. b\n. . c\n. . . d\n. c\n. . d\n";
  EXPECT_EQ(ex, ss.str());
}
