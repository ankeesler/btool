#ifndef BTOOL_NODE_STORE_H_
#define BTOOL_NODE_STORE_H_

#include <map>
#include <string>

#include "node/node.h"

namespace btool::node {

class Store {
 public:
  ~Store();

  Node *Create(const char *name);
  Node *Get(const char *name) const;

 private:
  std::map<std::string, Node *> nodes_;
};

};  // namespace btool::node

#endif  // BTOOL_NODE_STORE_H_
