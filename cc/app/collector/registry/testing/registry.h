#ifndef BTOOL_APP_COLLECTOR_REGISTRY_TESTING_REGISTRY_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_TESTING_REGISTRY_H_

#include "gmock/gmock.h"

#include "app/collector/registry/registry.h"
#include "app/collector/registry/serializer.h"

namespace btool::app::collector::registry::testing {

class MockRegistry : public ::btool::app::collector::registry::Registry {
 public:
  MOCK_METHOD1(GetIndex, void(Index *));
  MOCK_METHOD2(GetGaggle, void(std::string, Gaggle *));
};

class MockSerializer : public ::btool::app::collector::registry::Serializer {
 public:
  MOCK_METHOD2(UnmarshalIndex, void(std::istream *, Index *));
  MOCK_METHOD2(UnmarshalGaggle, void(std::istream *, Gaggle *));
};

};  // namespace btool::app::collector::registry::testing

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_TESTING_REGISTRY_H_
