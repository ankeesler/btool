#include "store.h"

#include <map>
#include <string>

#include "node/node.h"

namespace btool::app::collector {

Store::~Store() {
  for (auto kv : nodes_) {
    delete kv.second;
  }
}

::btool::node::Node *Store::Put(std::string name) {
  // TODO: this performance is bad?
  auto node = nodes_[name];
  if (node == nullptr) {
    node = new ::btool::node::Node(name);
    nodes_[node->name()] = node;
  }
  return node;
}

::btool::node::Node *Store::Get(const std::string &name) const {
  // TODO: this performance is bad?
  auto it = nodes_.find(name);
  if (it == nodes_.end()) {
    return nullptr;
  } else {
    return it->second;
  }
}

std::ostream &operator<<(std::ostream &os, const Store &s) {
  os << "[ ";
  for (auto kv : s.nodes_) {
    os << '\'' << kv.second->name() << '\'' << " ";
  }
  os << "]";
  return os;
}

};  // namespace btool::app::collector
