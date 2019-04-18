#include <iostream>

#include "gtest/gtest.h"
#include "gmock/gmock.h"
#include "error.h"
#include "cli/cli.h"
#include "cli/command.h"

using ::testing::Return;
using ::testing::ReturnRef;
using ::testing::StrictMock;

class MockCommand : public ::btool::cli::Command {
 public:
  MOCK_CONST_METHOD0(Name, const std::string&());
  MOCK_METHOD0(Run, ::btool::Error());
};

TEST(CLITest, Success) {
  ::btool::cli::CLI cli;

  StrictMock<MockCommand> first_command;
  const std::string first_command_name = "first-command";
  EXPECT_CALL(first_command, Name()).WillOnce(ReturnRef(first_command_name));
  cli.AddCommand(&first_command);

  StrictMock<MockCommand> second_command;
  const std::string second_command_name = "second-command";
  EXPECT_CALL(second_command, Name()).WillOnce(ReturnRef(second_command_name));
  EXPECT_CALL(second_command, Run()).WillOnce(Return(::btool::Error::Success()));
  cli.AddCommand(&second_command);

  const char *argv[] = {
    "--some-flag",
    "some-flag-arg",
    "--some-other-flag-arg",
    "some-other-flag-arg",
    "second-command",
    "--some-command-flag",
    "some-command-flag-arg",
    "--some-other-command-flag-arg",
    "some-other-command-flag-arg",
  };
  int argc = sizeof(argv)/sizeof(argv[0]);
  EXPECT_EQ(cli.Run(argc, argv), ::btool::Error::Success());
}
