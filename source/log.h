#ifndef BTOOL_LOG_H_
#define BTOOL_LOG_H_

#include <string>

namespace btool {

class Log {
 public:

  Log(const std::string& section): section_(section) { }

  void Println(const std::string& message);

 private:
  std::string section_;
};

}; // namespace btool

#endif // BTOOL_LOG_H_
