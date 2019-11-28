#ifndef BTOOL_APP_COLLECTOR_STORE_H_
#define BTOOL_APP_COLLECTOR_STORE_H_

#include <map>
#include <string>
#include <vector>

#include "node/node.h"

namespace btool::app::collector {

// Store
//
// Store is a bunch of Node's.
//
// The most important thing that a Store does is to delete Node heap data in
// its destructor. New nodes are created with Put(name).
class Store {
 public:
  ~Store();

  friend std::ostream &operator<<(std::ostream &os, const Store &s);

  // Put is idempotent: if no Node with the provided name exists, it will
  // create it; otherwise, the Node with the provided name will be returned.
  ::btool::node::Node *Put(std::string name);

  // Get returns nullptr iff no Node exists with the provided name.
  ::btool::node::Node *Get(const std::string &name) const;

  bool IsEmpty() const { return nodes_.empty(); }
  std::size_t Size() const { return nodes_.size(); }

  std::map<std::string, ::btool::node::Node *>::iterator begin() {
    return nodes_.begin();
  }

  std::map<std::string, ::btool::node::Node *>::iterator end() {
    return nodes_.end();
  }

 private:
  std::map<std::string, ::btool::node::Node *> nodes_;
};

std::ostream &operator<<(std::ostream &os, const Store &s);

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_STORE_H_
