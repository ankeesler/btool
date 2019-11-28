#ifndef BTOOL_APP_COLLECTOR_REGISTRY_SERIALIZER_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_SERIALIZER_H_

#include <istream>
#include <ostream>

#include "app/collector/registry/registry.h"

namespace btool::app::collector::registry {

class Serializer {
 public:
  virtual void UnmarshalIndex(std::istream *is, Index *i) = 0;
  virtual void UnmarshalGaggle(std::istream *is, Gaggle *g) = 0;
};

}  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_SERIALIZER_H_
