#include "lister.h"

#include <sstream>

// workaround for bug-00
#include "gmock/gmock.h"
#include "gtest/gtest.h"

// workaround for bug-00
#include "core/err.h"
#include "node/node.h"
#include "node/testing/node.h"

class ListerTest : public ::btool::node::testing::NodeTest {};

TEST_F(ListerTest, List) {
  std::stringstream ss;
  ::btool::app::lister::Lister l(&ss);

  l.List(a_);

  std::string ex = "a\n. b\n. . c\n. . . d\n. c\n. . d\n";
  EXPECT_EQ(ex, ss.str());
}
