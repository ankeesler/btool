#include "app/builder/work_pool_impl.h"

// bug-00
#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "err.h"
#include "node/node.h"
#include "node/testing/node.h"

using ::testing::HasSubstr;
using ::testing::UnorderedElementsAre;

class WorkPoolImplTest : public ::btool::node::testing::NodeTest {};

TEST_F(WorkPoolImplTest, Success) {
  ::btool::app::builder::WorkPoolImpl wpi(2);
  wpi.Submit([&]() -> const ::btool::node::Node* { return &a_; });
  EXPECT_EQ(&a_, wpi.Receive(nullptr));

  std::vector<const ::btool::node::Node*> nodes;
  wpi.Submit([&]() -> const ::btool::node::Node* { return &b_; });
  wpi.Submit([&]() -> const ::btool::node::Node* { return &c_; });
  nodes.push_back(wpi.Receive(nullptr));

  wpi.Submit([&]() -> const ::btool::node::Node* { return &d_; });
  nodes.push_back(wpi.Receive(nullptr));
  nodes.push_back(wpi.Receive(nullptr));

  EXPECT_THAT(nodes, UnorderedElementsAre(&b_, &c_, &d_));
}

TEST_F(WorkPoolImplTest, Failure) {
  ::btool::app::builder::WorkPoolImpl wpi(2);
  wpi.Submit([&]() -> const ::btool::node::Node* {
    THROW_ERR("an error was thrown here!");
    return nullptr;
  });

  ::btool::Err err;
  EXPECT_EQ(nullptr, wpi.Receive(&err));
  EXPECT_THAT(std::string(err.what()), HasSubstr("an error was thrown here!"));
}
