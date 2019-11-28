#ifndef BTOOL_UTIL_UNZIP_H_
#define BTOOL_UTIL_UNZIP_H_

#include <string>

namespace btool::util {

void Unzip(const std::string &zipfile, const std::string &dest_dir);

};  // namespace btool::util

#endif  // BTOOL_UTIL_UNZIP_H_
