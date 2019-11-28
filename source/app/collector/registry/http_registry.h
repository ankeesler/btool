#ifndef BTOOL_APP_COLLECTOR_REGISTRY_HTTPREGISTRY_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_HTTPREGISTRY_H_

#include <string>

#include "app/collector/registry/registry.h"
#include "app/collector/registry/serializer.h"

namespace btool::app::collector::registry {

class HttpRegistry : public Registry {
 public:
  HttpRegistry(std::string url, Serializer *s) : url_(url), s_(s) {}

  void GetIndex(Index *i) override;
  void GetGaggle(std::string name, Gaggle *i) override;

 private:
  std::string url_;
  Serializer *s_;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_HTTPREGISTRY_H_
