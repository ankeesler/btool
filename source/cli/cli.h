#include <vector>

#include "error.h"
#include "cli/command.h"

namespace btool::cli {

class CLI {
public:
  void AddCommand(const Command& command) {
    commands_.push_back(command);
  }

  Error Run();

private:
  std::vector<Command> commands_;
};

}; // namespace btool::cli
