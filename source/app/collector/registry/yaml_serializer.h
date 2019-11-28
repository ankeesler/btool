#ifndef BTOOL_APP_COLLECTOR_REGISTRY_YAMLSERIALIZER_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_YAMLSERIALIZER_H_

#include "app/collector/registry/serializer.h"

namespace btool::app::collector::registry {

class YamlSerializer : public Serializer {
 public:
  void UnmarshalIndex(std::istream *is, Index *i) override;
  void UnmarshalGaggle(std::istream *is, Gaggle *g) override;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_YAMLSERIALIZER_H_
