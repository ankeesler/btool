#include "util/string/string.h"

#include <algorithm>
#include <cstring>
#include <string>

namespace btool::util::string {

bool HasPrefix(const std::string &s, const std::string &prefix) {
  return (s.compare(0, prefix.size(), prefix) == 0);
}

bool HasSuffix(const std::string &s, const std::string &suffix) {
  return (s.size() >= suffix.size() &&
          s.compare(s.size() - suffix.size(), std::string::npos, suffix) == 0);
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
