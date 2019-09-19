#ifndef BTOOL_NODE_NODE_H_
#define BTOOL_NODE_NODE_H_

#include <map>
#include <string>
#include <vector>

namespace btool::node {

class Node {
public:
  Node(const std::string& name) : name_(name) { }

  const std::string& Name() const { return name_; }

private:
  std::string name_;
};

}; // namespace btool::node

#endif // BTOOL_NODE_NODE_H_
