#include "node/store.h"

#include <map>

namespace btool::node {

Store::~Store() {
  for (auto kv : nodes_) {
    delete kv.second;
  }
}

Node *Store::Create(const char *name) {
  auto node = new Node(name);
  nodes_[name] = node;
  return node;
}

Node *Store::Get(const char *name) const {
  auto it = nodes_.find(name);
  if (it == nodes_.end()) {
    return nullptr;
  } else {
    return it->second;
  }
}

};  // namespace btool::node
