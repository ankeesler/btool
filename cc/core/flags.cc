#include "core/flags.h"

#include <string>

#include "log.h"

namespace btool::core {

bool Flags::Parse(int argc, const char *argv[], std::string *err) {
  for (auto kv : bools_) {
    *kv.second = false;
  }

  for (int i = 0; i < argc; ++i) {
    const char *name = argv[i] + 1;  // move past '-'
    bool *value = bools_[name];
    DEBUG("found flag %s, %s value\n", name,
          (value == nullptr ? "NOT" : "also"));
    if (value != nullptr) {
      *value = true;
    }
  }
  return true;
}

};  // namespace btool::core
