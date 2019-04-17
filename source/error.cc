#include "error.h"

#include <string>

namespace btool {

Error Error::Create(const char *message) {
  Error e;
  e.exists_ = true;
  e.message_ = message;
  return e;
}

Error Error::Success() {
  Error e;
  e.exists_ = false;
  return e;
}

};
