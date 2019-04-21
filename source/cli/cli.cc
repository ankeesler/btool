#include "cli.h"

#include <cstring>
#include <vector>

#include "error.h"

namespace btool {
namespace cli {

static bool is_flag(const char *arg) {
  return ::strlen(arg) > 0 && arg[0] == '-';
}

Error CLI::Run(int argc, const char *argv[]) {
  Command *command = nullptr;
  std::vector<const char *> args;
  for (int i = 0; i < argc; i++) {
    if (is_flag(argv[i])) {
      log_->Debugf("found flag '%s' with arg '%s'\n", argv[i], argv[i+1]);
      i++;
      // TODO: process flag
    } else {
      log_->Debugf("command or arg: '%s'?\n", argv[i]);
      if (command == nullptr && (command = FindCommand(argv[i])) != nullptr) {
        log_->Debugln("command");
      } else {
        args.push_back(argv[i]);
        log_->Debugln("arg");
      }
    }
  }

  return (command != nullptr
          ? command->Run(args)
          : Error::Create("must supply sub command"));
}

Command *CLI::FindCommand(const char *arg) const {
  for (Command *command : commands_) {
    if (::strcmp(command->Name().c_str(), arg) == 0) {
      return command;
    }
  }
  return nullptr;
}

}; // namespace cli
}; // namespace btool
