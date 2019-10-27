#ifndef BTOOL_NODE_TESTING_NODE_H_
#define BTOOL_NODE_TESTING_NODE_H_

#include <memory>
#include <string>
#include <vector>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "core/err.h"
#include "node/node.h"

namespace btool::node::testing {

class MockResolver : public ::btool::node::Node::Resolver {
 public:
  MOCK_METHOD1(Resolve, ::btool::core::VoidErr(const ::btool::node::Node &));
};

class NodeTest : public ::testing::Test {
 protected:
  void SetUp() override {
    ::testing::Test::SetUp();
    a_.dependencies()->push_back(&b_);
    a_.dependencies()->push_back(&c_);
    a_.SetResolver(&ar_);

    tmpa_.dependencies()->push_back(&tmpb_);
    tmpa_.dependencies()->push_back(&tmpc_);

    b_.dependencies()->push_back(&c_);
    b_.SetResolver(&br_);

    tmpb_.dependencies()->push_back(&tmpc_);

    c_.dependencies()->push_back(&d_);
    c_.SetResolver(&cr_);

    tmpc_.dependencies()->push_back(&tmpd_);

    d_.SetResolver(&dr_);
  }

  ::btool::node::Node a_{"a"};
  ::btool::node::Node tmpa_{"/tmp/a"};
  MockResolver ar_;
  ::btool::node::Node b_{"b"};
  ::btool::node::Node tmpb_{"/tmp/b"};
  MockResolver br_;
  ::btool::node::Node c_{"c"};
  ::btool::node::Node tmpc_{"/tmp/c"};
  MockResolver cr_;
  ::btool::node::Node d_{"d"};
  ::btool::node::Node tmpd_{"/tmp/d"};
  MockResolver dr_;
};

};  // namespace btool::node::testing

#endif  // BTOOL_NODE_TESTING_NODE_H_
