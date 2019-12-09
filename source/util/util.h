#ifndef BTOOL_UTIL_UTIL_H_
#define BTOOL_UTIL_UTIL_H_

#include <algorithm>
#include <chrono>
#include <string>
#include <vector>

namespace btool::util {

template <typename T>
bool Contains(const std::vector<T> &v, T t) {
  return std::find(v.begin(), v.end(), t) != v.end();
}

std::string Hex(unsigned char *data, std::size_t data_size);

std::string CommaSeparatedNumber(std::size_t n);

};  // namespace btool::util

#endif  // BTOOL_UTIL_UTIL_H_
