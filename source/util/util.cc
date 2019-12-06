#include "util/util.h"

#include <cstdio>

#include <string>

namespace btool::util {

std::string Hex(const std::string &s) {
  std::string hex(s.size() * 2, '\0');
  for (std::size_t i = 0; i < s.size(); ++i) {
    ::sprintf(hex.data() + (i * 2), "%x", s[i]);
  }
  return hex;
}

};  // namespace btool::util
