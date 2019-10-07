#ifndef BTOOL_NODE_TESTING_NODE_H_
#define BTOOL_NODE_TESTING_NODE_H_

#include <string>
#include <vector>

#include "gmock/gmock.h"

#include "node/node.h"
#include "node/store.h"

namespace btool::node::testing {

class MockResolver : public ::btool::node::Node::Resolver {
 public:
  MOCK_METHOD2(Resolve, bool(const ::btool::node::Node &, std::string *));
};

std::unique_ptr<::btool::node::Store> Nodes0123();

};  // namespace btool::node::testing

#endif  // BTOOL_NODE_TESTING_NODE_H_
