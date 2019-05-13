#ifndef BTOOL_CLI_CLI_H_
#define BTOOL_CLI_CLI_H_

#include <vector>

#include "error.h"
#include "log.h"
#include "cli/command.h"

namespace btool {
namespace cli {

class CLI {
public:
  CLI(): log_(new Log("cli")) { }
  ~CLI() { delete log_; }

  void AddCommand(Command *command) {
    commands_.push_back(command);
  }

  Error Run(int argc, const char *argv[]);

private:
  btool::Log *log_;
  std::vector<Command *> commands_;

  Command *FindCommand(const char *arg) const;
};

}; // namespace cli
}; // namespace btool

#endif // BTOOL_CLI_CLI_H_
