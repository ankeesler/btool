#include "util/string/string.h"

#include <cstring>

namespace btool::util::string {

bool HasSuffix(const char *s, const char *suffix) {
  size_t s_len = ::strlen(s);
  size_t suffix_len = ::strlen(suffix);
  return (s_len >= suffix_len &&
          (::memcmp(s + s_len - suffix_len, suffix, suffix_len) == 0));
}

};  // namespace btool::util::string
