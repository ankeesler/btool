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

  l.List(*nodes->at(0));

  std::string ex = "a\n. b\n. . c\n. . . d\n. c\n. . d\n";
  EXPECT_EQ(ex, ss.str());

  ::btool::node::testing::Free(nodes);
}
