#ifndef LOG_H_
#define LOG_H_

#include <string>

class Log {
 public:

  Log(const std::string& section): section_(section) { }

  void Println(const std::string& message);

 private:
  std::string section_;
};

#endif // LOG_H_
