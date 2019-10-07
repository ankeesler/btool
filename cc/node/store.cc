#include "node/store.h"

#include <map>
#include <string>

namespace btool::node {

Store::~Store() {
  for (auto kv : nodes_) {
    delete kv.second;
  }
}

Node *Store::Create(const char *name) {
  // TODO: this performance is bad?
  std::string key(name);
  auto node = new Node(key);
  nodes_[name] = node;
  return node;
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

};  // namespace btool::node
