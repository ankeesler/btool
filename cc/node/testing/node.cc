#include "node.h"

#include <memory>

// workaround for bug-00
#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "node/node.h"
#include "node/store.h"

namespace btool::node::testing {

std::unique_ptr<::btool::node::Store> Nodes0123() {
  auto s = std::unique_ptr<::btool::node::Store>(new ::btool::node::Store);

  // 0 -> 1, 2
  // 1 -> 2
  // 2 -> 3
  // 4
  auto n3 = s->Create("3");
  auto n2 = s->Create("2");
  n2->AddDep(n3);
  auto n1 = s->Create("1");
  n1->AddDep(n2);
  auto n0 = s->Create("0");
  n0->AddDep(n1);
  n0->AddDep(n2);

  return s;
}

};  // namespace btool::node::testing
