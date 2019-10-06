#ifndef BTOOL_NODE_NODE_H_
#define BTOOL_NODE_NODE_H_

#include <map>
#include <ostream>
#include <set>
#include <string>
#include <vector>

namespace btool::node {

class Node {
 public:
  Node(const std::string &name) : name_(name) {}

  const std::string &Name() const { return name_; }
  void String(std::ostream *os) const;
  void Visit(std::function<void(const Node *)>) const;
  const std::vector<Node *> &Deps() const { return deps_; }

  void AddDep(Node *dep) { deps_.push_back(dep); }

 private:
  void String(std::ostream *os, int indent) const;
  void Visit(std::function<void(const Node *)>, std::set<const Node *> *) const;

  std::string name_;
  std::vector<Node *> deps_;
};

};  // namespace btool::node

#endif  // BTOOL_NODE_NODE_H_
