#ifndef BTOOL_UTIL_SHA256_H_
#define BTOOL_UTIL_SHA256_H_

#include <istream>
#include <string>

namespace btool::util {

std::string SHA256(std::istream *);

};  // namespace btool::util

#endif  // BTOOL_UTIL_SHA256_H_
