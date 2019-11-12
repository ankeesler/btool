#ifndef BTOOL_UTIL_FLAGS_H_
#define BTOOL_UTIL_FLAGS_H_

#include <cstring>

#include <map>
#include <string>

namespace btool::util {

class Flags {
 public:
  void Bool(std::string name, bool *value) { bools_[name] = value; }
  void String(std::string name, std::string *value) { strings_[name] = value; }

  // Returns true on success.
  bool Parse(int argc, const char *argv[], std::string *err);

 private:
  std::map<std::string, bool *> bools_;
  std::map<std::string, std::string *> strings_;
};

};  // namespace btool::util

#endif  // BTOOL_UTIL_FLAGS_H_
