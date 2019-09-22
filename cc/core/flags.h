#ifndef BTOOL_CORE_FLAGS_H_
#define BTOOL_CORE_FLAGS_H_

#include <map>
#include <string>

namespace btool::core {

class Flags {
 public:
  void Bool(const char *name, bool *value) { bools_[name] = value; }

  // Returns true on success.
  bool Parse(int argc, const char *argv[], std::string *err);

 private:
  struct cmp_str {
    bool operator()(const char *a, const char *b) const {
      return std::strcmp(a, b) < 0;
    }
  };
  std::map<const char *, bool *, cmp_str> bools_;
};

};  // namespace btool::core

#endif  // BTOOL_CORE_FLAGS_H_
