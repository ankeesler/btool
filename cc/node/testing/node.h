#ifndef BTOOL_NODE_TESTING_NODE_H_
#define BTOOL_NODE_TESTING_NODE_H_

#include <vector>

#include "node/node.h"

namespace btool::node::testing {

std::vector<::btool::node::Node *> *Nodes0123() {
  // 0 -> 1, 2
  // 1 -> 2
  // 2 -> 3
  // 4
  ::btool::node::Node *n3 = new ::btool::node::Node("d");
  ::btool::node::Node *n2 = new ::btool::node::Node("c");
  n2->AddDep(n3);
  ::btool::node::Node *n1 = new ::btool::node::Node("b");
  n1->AddDep(n2);
  ::btool::node::Node *n0 = new ::btool::node::Node("a");
  n0->AddDep(n1);
  n0->AddDep(n2);

  auto nodes = new std::vector<::btool::node::Node *>;
  nodes->push_back(n0);
  nodes->push_back(n1);
  nodes->push_back(n2);
  nodes->push_back(n3);

  return nodes;
}

void Free(std::vector<::btool::node::Node *> *nodes) {
  for (auto node : *nodes) {
    delete node;
  }
  delete nodes;
}

};  // namespace btool::node::testing

#endif  // BTOOL_NODE_TESTING_NODE_H_
