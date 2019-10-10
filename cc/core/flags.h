#ifndef BTOOL_CORE_FLAGS_H_
#define BTOOL_CORE_FLAGS_H_

#include <cstring>

#include <map>
#include <string>

namespace btool::core {

class Flags {
 public:
  void Bool(const char *name, bool *value) { bools_[name] = value; }
  void String(const char *name, std::string *value) { strings_[name] = value; }

  // Returns true on success.
  bool Parse(int argc, const char *argv[], std::string *err);

 private:
  struct cmp_str {
    bool operator()(const char *a, const char *b) const {
      return ::strcmp(a, b) < 0;
    }
  };
  std::map<const char *, bool *, cmp_str> bools_;
  std::map<const char *, std::string *, cmp_str> strings_;
};

};  // namespace btool::core

#endif  // BTOOL_CORE_FLAGS_H_
