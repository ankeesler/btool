#ifndef BTOOL_APP_COLLECTOR_REGISTRY_YAMLSERIALIZER_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_YAMLSERIALIZER_H_

#include <istream>
#include <ostream>

#include "app/collector/registry/serializer.h"
#include "err.h"

namespace btool::app::collector::registry {

template <typename T>
class YamlSerializer : public Serializer<T> {
 public:
  void Unmarshal(std::istream *is, T *t) override;
  void Marshal(std::ostream *os, const T &t) override;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_YAMLSERIALIZER_H_
