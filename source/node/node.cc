#include "node/node.h"

#include <functional>
#include <ostream>
#include <set>

namespace btool::node {

void Node::String(std::ostream *os) const { String(os, 0); }

void Node::String(std::ostream *os, int indent) const {
  for (int i = 0; i < indent; ++i) {
    *os << ". ";
  }
  *os << name_ << std::endl;
  for (auto dep : dependencies_) {
    dep->String(os, indent + 1);
  }
}

void Node::Visit(std::function<void(const Node *)> f) const {
  std::set<const Node *> visited;
  Visit(f, &visited);
}

void Node::Visit(std::function<void(const Node *)> f,
                 std::set<const Node *> *visited) const {
  if (visited->count(this) == 0) {
    for (auto dep : dependencies_) {
      dep->Visit(f, visited);
    }
    f(this);
  }
  visited->insert(this);
}

};  // namespace btool::node
