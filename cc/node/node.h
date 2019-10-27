#ifndef BTOOL_NODE_NODE_H_
#define BTOOL_NODE_NODE_H_

#include <functional>
#include <map>
#include <ostream>
#include <set>
#include <string>
#include <vector>

#include "core/err.h"
#include "node/property_store.h"

namespace btool::node {

class Node {
 public:
  class Resolver {
   public:
    ~Resolver() {}
    virtual ::btool::core::VoidErr Resolve(const Node &node) = 0;
  };

  Node(std::string name) : name_(name), resolver_(nullptr) {}

  const std::string &Name() const { return name_; }
  void String(std::ostream *os) const;
  void Visit(std::function<void(const Node *)>) const;
  Resolver *resolver() const { return resolver_; }

  std::vector<Node *> *dependencies() { return &dependencies_; }
  PropertyStore *property_store() { return &property_store_; }

  void SetResolver(Resolver *resolver) { resolver_ = resolver; }

 private:
  void String(std::ostream *os, int indent) const;
  void Visit(std::function<void(const Node *)>, std::set<const Node *> *) const;

  std::string name_;
  std::vector<Node *> dependencies_;
  Resolver *resolver_;
  PropertyStore property_store_;
};

};  // namespace btool::node

#endif  // BTOOL_NODE_NODE_H_
