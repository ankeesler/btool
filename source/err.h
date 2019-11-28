#ifndef BTOOL_ERR_H_
#define BTOOL_ERR_H_

#include <iostream>
#include <optional>
#include <ostream>
#include <stdexcept>
#include <string>

namespace btool {

class Err : public std::runtime_error {
 public:
  Err(std::string what, std::string file, int line)
      : std::runtime_error(what + " @ " + file + ":" + std::to_string(line)) {}
};

#define THROW_ERR(what) throw Err((what), __FILE__, __LINE__)

};  // namespace btool

#endif  // BTOOL_ERR_H_
