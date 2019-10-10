#include "err.h"

#include <ostream>

namespace btool::core {

std::ostream &operator<<(std::ostream &os, const VoidErr &err) {
  if (err) {
    os << "err: failure: " << err.Msg();
  } else {
    os << "err: success";
  }
  return os;
}

};  // namespace btool::core
