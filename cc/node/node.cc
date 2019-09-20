#include "node/node.h"

#include <ostream>

namespace btool::node {

void Node::String(std::ostream *os) const {
  String(os, 0);
}

void Node::String(std::ostream *os, int indent) const {
  for (int i = 0; i < indent; ++i) {
    *os << ". ";
  }
  *os << name_ << std::endl;
  for (auto dep : deps_) {
    dep->String(os, indent + 1);
  }
}

}; // namespace btool::node
