#ifndef BTOOL_LOG_H_
#define BTOOL_LOG_H_

#include <cstdarg>
#include <string>

namespace btool {

class Log {
 public:

  enum Level {
    DEBUG = 0,
    INFO = 1,
  };

  static Level level;

  Log(const std::string& section): section_(section) { }

  void Debugln(const std::string& message) {
    Println(DEBUG, message);
  }

  void Debug(const std::string& message) {
    Print(DEBUG, message);
  }

  void Debugf(const char *format, ...);

 private:
  std::string section_;

  void Println(Level level, const std::string& message);
  void Print(Level level, const std::string& message);
};

}; // namespace btool

#endif // BTOOL_LOG_H_
