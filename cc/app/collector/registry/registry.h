#ifndef BTOOL_APP_COLLECTOR_REGISTRY_REGISTRY_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_REGISTRY_H_

#include <vector>

#include "node/property_store.h"

namespace btool::app::collector::registry {

struct IndexFile {
  std::string path;
  std::string sha256;

  bool operator==(const IndexFile &if0) const {
    return if0.path == path && if0.sha256 == sha256;
  }
};

struct Index {
  std::vector<IndexFile> files;

  bool operator==(const Index &i0) const { return i0.files == files; }
};

struct Resolver {
  std::string name;
  ::btool::node::PropertyStore config;
};

struct Node {
  std::string name;
  std::vector<std::string> dependencies;
  ::btool::node::PropertyStore property_store;
  Resolver resolver;
};

struct Gaggle {
  std::vector<Node> nodes;
};

class Registry {
 public:
  virtual void GetIndex(Index *i) = 0;
  virtual void GetGaggle(std::string path, Gaggle *g) = 0;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_REGISTRY_H_
