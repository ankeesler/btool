#ifndef BTOOL_CLI_COMMAND_H_
#define BTOOL_CLI_COMMAND_H_

#include <string>

#include "error.h"

namespace btool {
namespace cli {

class Command {
public:
  virtual ~Command() { }
  virtual const ::std::string& Name() const = 0;
  virtual ::btool::Error Run() = 0;
};

}; // namespace cli
}; // namespace btool

#endif // BTOOL_CLI_COMMAND_H_
