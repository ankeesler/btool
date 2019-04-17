#include "log.h"

#include <iostream>

void Log::Println(const std::string& message) {
  std::cout << "[" << section_ << "] " << message << std::endl;
}
