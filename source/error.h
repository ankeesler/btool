#include <string>

namespace btool {

class Error {
 public:
  static Error Create(const char *message);
  static Error Success();

  bool Exists() const { return exists_; }
  const std::string& Message() const { return message_; }

 private:
  bool exists_;
  std::string message_;
};

}; // namespace btool
