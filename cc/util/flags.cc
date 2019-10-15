#include "flags.h"

#include <string>

#include "log.h"

namespace btool::util {

bool Flags::Parse(int argc, const char *argv[], std::string *err) {
  for (auto kv : bools_) {
    *kv.second = false;
  }

  for (int i = 0; i < argc; ++i) {
    const char *name = argv[i] + 1;  // move past '-'
    bool *bv = bools_[name];
    DEBUG("found bool flag %s, %s value\n", name,
          (bv == nullptr ? "NOT" : "also"));
    if (bv != nullptr) {
      *bv = true;
    }

    std::string *sv = strings_[name];
    DEBUG("found string flag %s, %s value\n", name,
          (sv == nullptr ? "NOT" : "also"));
    if (sv != nullptr) {
      *sv = argv[++i];
    }
  }
  return true;
}

};  // namespace btool::util
