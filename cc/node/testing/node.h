#ifndef BTOOL_NODE_TESTING_NODE_H_
#define BTOOL_NODE_TESTING_NODE_H_

#include <memory>
#include <string>
#include <vector>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "core/err.h"
#include "node/node.h"
#include "node/store.h"

namespace btool::node::testing {

class MockResolver : public ::btool::node::Node::Resolver {
 public:
  MOCK_METHOD1(Resolve, ::btool::core::VoidErr(const ::btool::node::Node &));
};

class NodeTest : public ::testing::Test {
 protected:
  void SetUp() override {
    ::testing::Test::SetUp();
    a_.AddDep(&b_);
    a_.AddDep(&c_);
    a_.SetResolver(&ar_);

    b_.AddDep(&c_);
    b_.SetResolver(&br_);

    c_.AddDep(&d_);
    c_.SetResolver(&cr_);

    d_.SetResolver(&dr_);
  }

  ::btool::node::Node a_{"a"};
  MockResolver ar_;
  ::btool::node::Node b_{"b"};
  MockResolver br_;
  ::btool::node::Node c_{"c"};
  MockResolver cr_;
  ::btool::node::Node d_{"d"};
  MockResolver dr_;
};

};  // namespace btool::node::testing

#endif  // BTOOL_NODE_TESTING_NODE_H_
