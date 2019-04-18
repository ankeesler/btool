#include "cli.h"

#include <cstring>

#include "error.h"

namespace btool {
namespace cli {

static bool is_flag(const char *arg) {
  return ::strlen(arg) > 0 && arg[0] == '-';
}

Error CLI::Run(int argc, const char *argv[]) {
  for (int i = 0; i < argc; i++) {
    if (is_flag(argv[i])) {
      log_->Println("found flag:");
      log_->Println(argv[i]);
      i++;
      log_->Println("with arg:");
      log_->Println(argv[i]);
      // TODO: process flag
    } else {
      log_->Println("searching for command:");
      log_->Println(argv[i]);

      Command *command = FindCommand(argv[i]);
      if (command != NULL) {
        return command->Run();
      }

      log_->Println("didn't find it");
    }
  }

  return Error::Create("must supply sub command");
}

Command *CLI::FindCommand(const char *arg) const {
  for (Command *command : commands_) {
    if (::strcmp(command->Name().c_str(), arg) == 0) {
      return command;
    }
  }
  return NULL;
}

}; // namespace cli
}; // namespace btool
