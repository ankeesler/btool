#ifndef BTOOL_ERROR_H_
#define BTOOL_ERROR_H_

#include <string>

namespace btool {

class Error {
 public:
  static Error Create(const char *message);
  static Error Success();

  bool Exists() const { return exists_; }
  const std::string& Message() const { return message_; }

  bool operator==(const Error& e) const {
    return exists_ == e.exists_ && message_ == e.message_;
  }

  bool operator!=(const Error& e) const {
    return exists_ != e.exists_ || message_ != e.message_;
  }

 private:
  bool exists_;
  std::string message_;
};

}; // namespace btool

#endif // BTOOL_ERROR_H_
