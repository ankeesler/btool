#ifndef BTOOL_CREATE_CLASS_COMMAND_H_
#define BTOOL_CREATE_CLASS_COMMAND_H_

#include <string>
#include <vector>

#include "cli/command.h"
#include "fs.h"

namespace btool {

class CreateClassCommand : public btool::cli::Command {
public:
  CreateClassCommand(FS *fs): fs_(fs) { }

  const std::string& Name() const {
    return name;
  }
  
  Error Run(const std::vector<const char *>& args);

private:
  static const std::string name;

  FS *fs_;
};

}; // namespace btool

#endif // BTOOL_CREATE_CLASS_COMMAND_H_
