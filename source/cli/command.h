#ifndef BTOOL_CLI_COMMAND_H_
#define BTOOL_CLI_COMMAND_H_

#include <string>
#include <vector>

#include "error.h"

namespace btool {
namespace cli {

class Command {
public:
  virtual ~Command() { }
  virtual const ::std::string& Name() const = 0;
  virtual btool::Error Run(const std::vector<const char *>&) = 0;
};

}; // namespace cli
}; // namespace btool

#endif // BTOOL_CLI_COMMAND_H_
