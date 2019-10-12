#include "node/store.h"

#include <map>
#include <string>

namespace btool::node {

Store::~Store() {
  for (auto kv : nodes_) {
    delete kv.second;
  }
}

Node *Store::Put(const char *name) {
  // TODO: this performance is bad?
  std::string key(name);
  auto node = nodes_[key];
  if (node == nullptr) {
    node = new Node(name);
    Set(node);
  }
  return node;
}

void Store::Set(Node *node) {
  nodes_[node->Name()] = node;
  for (auto l : ls_) {
    l->OnSet(this, node->Name());
  }
}

Node *Store::Get(const char *name) const {
  // TODO: this performance is bad?
  std::string key(name);
  auto it = nodes_.find(key);
  if (it == nodes_.end()) {
    return nullptr;
  } else {
    return it->second;
  }
}

std::ostream &operator<<(std::ostream &os, const Store &s) {
  os << "[ ";
  for (auto kv : s.nodes_) {
    os << '\'' << kv.second->Name() << '\'' << " ";
  }
  os << "]";
  return os;
}

};  // namespace btool::node