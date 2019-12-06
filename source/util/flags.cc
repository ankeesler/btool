#include "util/flags.h"

#include <iostream>
#include <ostream>
#include <string>

#include "log.h"

namespace btool::util {

bool Flags::Parse(int argc, const char *argv[], std::string *err) {
  for (auto f : bool_flags_) {
    *f.value_ = false;
  }

  for (int i = 0; i < argc; ++i) {
    const char *name = argv[i] + 1;  // move past '-'

    Flags::Flag<bool> *fb = Find(&bool_flags_, name);
    if (fb != nullptr) {
      DEBUGS() << "found bool flag named " << name << std::endl;
      *fb->value_ = true;
      continue;
    }

    Flags::Flag<std::string> *fs = Find(&string_flags_, name);
    if (fs != nullptr) {
      const char *value = argv[++i];
      DEBUGS() << "found string flag named " << name << " with value " << value
               << std::endl;
      *fs->value_ = value;
      continue;
    }

    Flags::Flag<int> *fi = Find(&int_flags_, name);
    if (fi != nullptr) {
      const char *value = argv[++i];
      DEBUGS() << "found int flag named " << name << " with value " << value
               << std::endl;
      *fi->value_ = std::stoi(value);
      continue;
    }
  }
  return true;
}

void Flags::Usage(std::ostream *os) {
  for (const auto &f : bool_flags_) {
    Usage(os, f);
  }

  for (const auto &f : string_flags_) {
    Usage(os, f);
  }

  for (const auto &f : int_flags_) {
    Usage(os, f);
  }
}

};  // namespace btool::util
