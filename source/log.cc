#include "log.h"

#include <cstdarg>
#include <iostream>
#include <string>

namespace btool {

void Log::Println(const std::string& message) {
  std::cout << "[" << section_ << "] " << message << std::endl;
}

void Log::Print(const std::string& message) {
  std::cout << "[" << section_ << "] " << message;
}

void Log::Printf(const std::string& format, ...) {
  std::cout << "[" << section_ << "] ";

  va_list args;
  va_start(args, format);
  vprintf(format.c_str(), args);
  va_end(args);
}

}; // namespace btool
