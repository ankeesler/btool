#ifndef BTOOL_UTIL_TESTING_UTIL_H_
#define BTOOL_UTIL_TESTING_UTIL_H_

#include <string>

#include "gmock/gmock.h"

#include "util/cache.h"

namespace btool::util::testing {

template <typename T>
class MockCache : public ::btool::util::Cache<T> {
 public:
  MOCK_METHOD2_T(Get, bool(const std::string &, T *));
  MOCK_METHOD2_T(Set, void(const std::string &, const T &));
};

};  // namespace btool::util::testing

#endif  // BTOOL_UTIL_TESTING_UTIL_H_
