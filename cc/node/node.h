#ifndef BTOOL_NODE_NODE_H_
#define BTOOL_NODE_NODE_H_

#include <functional>
#include <map>
#include <ostream>
#include <set>
#include <string>
#include <vector>

#include "node/property_store.h"

namespace btool::node {

class Node {
 public:
  class Resolver {
   public:
    virtual ~Resolver() {}
    virtual void Resolve(const Node &node) = 0;
  };

  Node(std::string name) : name_(name), resolver_(nullptr) {}

  const std::string &name() const { return name_; }
  std::vector<Node *> *dependencies() { return &dependencies_; }
  const std::vector<Node *> *dependencies() const { return &dependencies_; }
  Resolver *resolver() const { return resolver_; }
  void set_resolver(Resolver *resolver) { resolver_ = resolver; }
  const PropertyStore *property_store() const { return &property_store_; }
  PropertyStore *property_store() { return &property_store_; }

  void String(std::ostream *os) const;
  void Visit(std::function<void(const Node *)>) const;

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
