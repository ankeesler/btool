#ifndef BTOOL_UTIL_STRING_STRING_H_
#define BTOOL_UTIL_STRING_STRING_H_

#include <string>

namespace btool::util::string {

bool HasSuffix(const char *s, const char *suffix);
std::string Replace(const std::string &s, const std::string &from,
                    const std::string &to);

};  // namespace btool::util::string

#endif  // BTOOL_UTIL_STRING_STRING_H_
