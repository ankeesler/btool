#include "lister.h"

#include <sstream>

// workaround for bug-00
#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "node/node.h"
#include "node/testing/node.h"

TEST(Lister, List) {
  std::stringstream ss;
  ::btool::app::lister::Lister l(&ss);

  auto nodes = ::btool::node::testing::Nodes0123();

  l.List(*nodes->Get("0"));

  std::string ex = "0\n. 1\n. . 2\n. . . 3\n. 2\n. . 3\n";
  EXPECT_EQ(ex, ss.str());
}
