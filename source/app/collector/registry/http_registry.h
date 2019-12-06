#ifndef BTOOL_APP_COLLECTOR_REGISTRY_HTTPREGISTRY_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_HTTPREGISTRY_H_

#include <string>

#include "app/collector/registry/registry.h"
#include "app/collector/registry/serializer.h"

namespace btool::app::collector::registry {

class HttpRegistry : public Registry {
 public:
  HttpRegistry(std::string url, Serializer<Index> *s_i, Serializer<Gaggle> *s_g)
      : url_(url), s_i_(s_i), s_g_(s_g) {}

  std::string GetName() override { return url_; }
  void GetIndex(Index *i) override;
  void GetGaggle(std::string name, Gaggle *i) override;

 private:
  std::string url_;
  Serializer<Index> *s_i_;
  Serializer<Gaggle> *s_g_;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_HTTPREGISTRY_H_
