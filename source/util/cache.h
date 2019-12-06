#ifndef BTOOL_UTIL_CACHE_H_
#define BTOOL_UTIL_CACHE_H_

#include <string>

namespace btool::util {

template <typename T>
class Cache {
 public:
  virtual bool Get(const std::string &name, T *t) = 0;
  virtual void Set(const std::string &name, const T &t) = 0;
};

};  // namespace btool::util

#endif  // BTOOL_UTIL_CACHE_H_
