#include "log.h"

#include <cstdarg>
#include <iostream>
#include <string>

namespace btool {

Log::Level Log::level = Log::Level::INFO;

void Log::Println(Log::Level level, const std::string& message) {
  if (Log::level <= level) {
    std::cout << "[" << section_ << "] "
              << "[" << level << "] "
              << message << std::endl;
  }
}

void Log::Print(Log::Level level, const std::string& message) {
  if (Log::level <= level) {
    std::cout << "[" << section_ << "] "
              << "[" << level << "] "
              << message;
  }
}

void Log::Debugf(const char *format, ...) {
  if (Log::level <= DEBUG) {
    std::cout << "[" << section_ << "] " << "[" << Log::Level::DEBUG << "] ";
    va_list args;
    va_start(args, format);
    vprintf(format, args);
    va_end(args);
  }
}

}; // namespace btool
