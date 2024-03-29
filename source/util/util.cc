#include "util/util.h"

#include <chrono>
#include <cstdio>
#include <functional>
#include <string>

namespace btool::util {

std::string Hex(unsigned char *data, std::size_t data_size) {
  std::string hex(data_size * 2, '\0');
  for (std::size_t i = 0; i < data_size; ++i) {
    ::sprintf(hex.data() + (i * 2), "%x%x", (data[i] & 0xF0) >> 4,
              data[i] & 0x0F);
  }
  return hex;
}

std::string CommaSeparatedNumber(std::size_t n) {
  std::string s = std::to_string(n);
  std::size_t size = s.size();
  while (size > 3) {
    s.insert(size - 3, ",");
    size -= 3;
  }
  return s;
}

std::chrono::steady_clock::duration Time(std::function<void()> f) {
  auto start = std::chrono::steady_clock::now();
  f();
  auto end = std::chrono::steady_clock::now();
  return end - start;
}

};  // namespace btool::util
