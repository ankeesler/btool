#ifndef BTOOL_CREATE_CLASS_COMMAND_H_
#define BTOOL_CREATE_CLASS_COMMAND_H_

#include <string>
#include <vector>

#include "cli/command.h"

namespace btool {

class CreateClassCommand : public btool::cli::Command {
public:
  virtual const std::string& Name() const {
    return name;
  }
  
  virtual btool::Error Run(const std::vector<const char *>& args);

private:
  static const std::string name;
};

}; // namespace btool

#endif // BTOOL_CREATE_CLASS_COMMAND_H_
