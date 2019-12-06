#ifndef BTOOL_APP_COLLECTOR_REGISTRY_SERIALIZER_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_SERIALIZER_H_

#include <istream>
#include <ostream>

#include "app/collector/registry/registry.h"

namespace btool::app::collector::registry {

template <typename T>
class Serializer {
 public:
  virtual void Unmarshal(std::istream *is, T *t) = 0;
  virtual void Marshal(std::ostream *is, const T &t) = 0;
};

}  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_SERIALIZER_H_
