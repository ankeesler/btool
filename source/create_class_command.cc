#include "create_class_command.h"

namespace btool {

const std::string CreateClassCommand::name = "create-class";

btool::Error CreateClassCommand::Run(const std::vector<const char *>& args) {
  return btool::Error::Success();
}

}; // namespace btool
