#ifndef BTOOL_APP_COLLECTOR_REGISTRY_FSREGISTRY_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_FSREGISTRY_H_

#include <map>
#include <string>

#include "app/collector/registry/registry.h"
#include "app/collector/registry/serializer.h"
#include "err.h"

namespace btool::app::collector::registry {

class FsRegistry : public Registry {
 public:
  FsRegistry(std::string root, Serializer *s)
      : root_(root), s_(s), initialized_(false) {}

  void GetIndex(Index *i) override {
    if (!initialized_) {
      Initialize();
    }

    *i = i_;
  }

  void GetGaggle(std::string path, Gaggle *g) override {
    if (!initialized_) {
      Initialize();
    }

    auto it = gs_.find(path);
    if (it == gs_.end()) {
      THROW_ERR("unknown gaggle for path " + path);
    } else {
      *g = it->second;
    }
  }

 private:
  void Initialize();

  std::string root_;
  Serializer *s_;

  bool initialized_;
  Index i_;
  std::map<std::string, Gaggle> gs_;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_FSREGISTRY_H_
