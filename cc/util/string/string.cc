#include "util/string/string.h"

#include <cstring>
#include <string>

namespace btool::util::string {

bool HasSuffix(const char *s, const char *suffix) {
  size_t s_len = ::strlen(s);
  size_t suffix_len = ::strlen(suffix);
  return (s_len >= suffix_len &&
          (::memcmp(s + s_len - suffix_len, suffix, suffix_len) == 0));
}

std::string Replace(const std::string &s, const std::string &from,
                    const std::string &to) {
  std::size_t index = s.find(from);
  if (index == std::string::npos) {
    return s;
  } else {
    std::string copy(s);
    copy.replace(index, std::string::npos, to);
    return copy;
  }
}

};  // namespace btool::util::string
