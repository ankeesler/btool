#include "log.h"

#include <iostream>
#include <string>

namespace btool {

void Log::Println(const std::string& message) {
  std::cout << "[" << section_ << "] " << message << std::endl;
}

}; // namespace btool
