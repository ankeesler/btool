#ifndef BTOOL_UTIL_FLAGS_H_
#define BTOOL_UTIL_FLAGS_H_

#include <cstring>

#include <iostream>
#include <map>
#include <ostream>
#include <string>
#include <vector>

namespace btool::util {

class Flags {
 public:
  void Bool(std::string name, std::string description, bool *value) {
    bool_flags_.push_back(Flag(name, description, value));
  }
  void String(std::string name, std::string description, std::string *value) {
    string_flags_.push_back(Flag(name, description, value));
  }
  void Int(std::string name, std::string description, int *value) {
    int_flags_.push_back(Flag(name, description, value));
  }

  // Returns true on success.
  bool Parse(int argc, const char *argv[], std::string *err);

  void Usage(std::ostream *os);

 private:
  template <typename T>
  class Flag {
   public:
    Flag(std::string name, std::string description, T *value)
        : name_(name), description_(description), value_(value) {}

    std::string name_;
    std::string description_;
    T *value_;
  };

  template <typename T>
  static Flag<T> *Find(std::vector<Flag<T>> *flags, const char *name) {
    for (std::size_t i = 0; i < flags->size(); ++i) {
      if (flags->at(i).name_ == name) {
        return flags->data() + i;
      }
    }
    return nullptr;
  }

  template <typename T>
  static void Usage(std::ostream *os, const Flag<T> &f) {
    *os << "  "
        << "-" << f.name_ << std::endl;
    *os << "  "
        << "  " << f.description_ << std::endl;
  }

  std::vector<Flag<bool>> bool_flags_;
  std::vector<Flag<std::string>> string_flags_;
  std::vector<Flag<int>> int_flags_;
};

};  // namespace btool::util

#endif  // BTOOL_UTIL_FLAGS_H_
