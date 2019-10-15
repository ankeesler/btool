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
//
// Store provides a Listener interface so that clients can be notified of Store
// events.
class Store {
 public:
  class Listener {
   public:
    ~Listener() {}

    // OnSet notifies a Listener that a Node with the provided name has been
    // Set() to the provided Store.
    //
    // Note that a call to Put() will trigger a Set() call, which will trigger a
    // call to OnSet().
    virtual void OnSet(Store *, const std::string &name) = 0;
  };

  ~Store();

  friend std::ostream &operator<<(std::ostream &os, const Store &s);

  // Put is idempotent: if no Node with the provided name exists, it will
  // create it; otherwise, the Node with the provided name will be returned.
  ::btool::node::Node *Put(std::string name);

  // Set sets the provided Node to the Store.
  void Set(::btool::node::Node *node);

  // Get returns nullptr iff no Node exists with the provided name.
  ::btool::node::Node *Get(const std::string &name) const;

  // Listen adds a Listener to this Store. The Listener will be notified when
  // about various Store events.
  void Listen(Listener *l) { ls_.push_back(l); }

  bool IsEmpty() const { return nodes_.empty(); }

 private:
  std::map<std::string, ::btool::node::Node *> nodes_;
  std::vector<Listener *> ls_;
};

std::ostream &operator<<(std::ostream &os, const Store &s);

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_STORE_H_
