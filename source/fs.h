#ifndef BTOOL_FS_H_
#define BTOOL_FS_H_

#include <string>

#include "error.h"

namespace btool {

class FS {
public:
  virtual Error WriteFile(const std::string& file, const std::string& contents) = 0;
};

}; // namespace btool

#endif // BTOOL_FS_H
