#ifndef BTOOL_APP_COLLECTOR_REGISTRY_TESTING_REGISTRY_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_TESTING_REGISTRY_H_

#include <istream>
#include <ostream>
#include <string>

#include "gmock/gmock.h"

#include "app/collector/registry/registry.h"
#include "app/collector/registry/serializer.h"

namespace btool::app::collector::registry::testing {

class MockRegistry : public ::btool::app::collector::registry::Registry {
 public:
  MOCK_METHOD0(GetName, std::string());
  MOCK_METHOD1(GetIndex, void(Index *));
  MOCK_METHOD2(GetGaggle, void(std::string, Gaggle *));
};

template <typename T>
class MockSerializer : public ::btool::app::collector::registry::Serializer<T> {
 public:
  MOCK_METHOD2_T(Unmarshal, void(std::istream *, T *));
  MOCK_METHOD2_T(Marshal, void(std::ostream *, const T &));
};

};  // namespace btool::app::collector::registry::testing

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_TESTING_REGISTRY_H_
