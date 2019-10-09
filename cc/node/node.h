#ifndef BTOOL_NODE_NODE_H_
#define BTOOL_NODE_NODE_H_

#include <functional>
#include <map>
#include <ostream>
#include <set>
#include <string>
#include <vector>

#include "core/err.h"

namespace btool::node {

class Node {
 public:
  class Resolver {
   public:
    ~Resolver() {}
    virtual ::btool::core::VoidErr Resolve(const Node &node) = 0;
  };

  Node(const char *name) : name_(name) {}

  const std::string &Name() const { return name_; }
  void String(std::ostream *os) const;
  void Visit(std::function<void(const Node *)>) const;
  const std::vector<Node *> &Deps() const { return deps_; }
  Resolver *Resolver() const { return resolver_; }

  void AddDep(Node *dep) { deps_.push_back(dep); }
  void SetResolver(class Resolver *resolver) { resolver_ = resolver; }

 private:
  void String(std::ostream *os, int indent) const;
  void Visit(std::function<void(const Node *)>, std::set<const Node *> *) const;

  std::string name_;
  std::vector<Node *> deps_;
  class Resolver *resolver_;
};

};  // namespace btool::node

#endif  // BTOOL_NODE_NODE_H_
