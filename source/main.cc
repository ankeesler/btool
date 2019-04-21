#include "cli/cli.h"
#include "create_class_command.h"
#include "log.h"

int main(int argc, const char *argv[]) {
  btool::Log log("main");
  log.Println("start");
  
  btool::CreateClassCommand create_class_command;

  btool::cli::CLI cli;
  cli.AddCommand(&create_class_command);
  cli.Run(argc - 1, argv + 1);

  log.Println("end");
  return 0;
}
