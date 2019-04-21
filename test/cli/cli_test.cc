#include <iostream>
#include <vector>

#include "gtest/gtest.h"
#include "gmock/gmock.h"
#include "error.h"
#include "cli/cli.h"
#include "cli/command.h"

using ::testing::Return;
using ::testing::ReturnRef;
using ::testing::StrictMock;

class MockCommand : public btool::cli::Command {
 public:
  MOCK_CONST_METHOD0(Name, const std::string&());
  MOCK_METHOD1(Run, btool::Error(const std::vector<const char *>&));
};

TEST(CLITest, Success) {
  btool::cli::CLI cli;

  StrictMock<MockCommand> first_command;
  const std::string first_command_name = "first-command";
  EXPECT_CALL(first_command, Name())
    .WillOnce(ReturnRef(first_command_name));
  cli.AddCommand(&first_command);

  StrictMock<MockCommand> second_command;
  const std::string second_command_name = "second-command";
  const std::vector<const char *> args = { "some-actual-arg" };
  EXPECT_CALL(second_command, Name())
    .WillOnce(ReturnRef(second_command_name));
  EXPECT_CALL(second_command, Run(args))
    .WillOnce(Return(btool::Error::Success()));
  cli.AddCommand(&second_command);

  const char *argv[] = {
    "--some-flag",
    "some-flag-arg",
    "--some-other-flag-arg",
    "some-other-flag-arg",
    "second-command",
    "--some-command-flag",
    "some-command-flag-arg",
    "some-actual-arg",
    "--some-other-command-flag-arg",
    "some-other-command-flag-arg",
  };
  int argc = sizeof(argv)/sizeof(argv[0]);
  EXPECT_EQ(cli.Run(argc, argv), btool::Error::Success());
}

TEST(CLITest, UnknownCommand) {
  btool::cli::CLI cli;

  StrictMock<MockCommand> first_command;
  const std::string first_command_name = "first-command";
  EXPECT_CALL(first_command, Name())
    .WillRepeatedly(ReturnRef(first_command_name));
  cli.AddCommand(&first_command);

  StrictMock<MockCommand> second_command;
  const std::string second_command_name = "second-command";
  EXPECT_CALL(second_command, Name())
    .WillRepeatedly(ReturnRef(second_command_name));
  cli.AddCommand(&second_command);

  const char *argv[] = {
    "--some-flag",
    "some-flag-arg",
    "--some-other-flag-arg",
    "some-other-flag-arg",
    "unknown-command",
    "--some-command-flag",
    "some-command-flag-arg",
    "some-actual-arg",
    "--some-other-command-flag-arg",
    "some-other-command-flag-arg",
  };
  int argc = sizeof(argv)/sizeof(argv[0]);
  EXPECT_TRUE(cli.Run(argc, argv).Exists());
}

TEST(CLITest, BadCommand) {
  btool::cli::CLI cli;

  StrictMock<MockCommand> first_command;
  const std::string first_command_name = "first-command";
  EXPECT_CALL(first_command, Name())
    .WillRepeatedly(ReturnRef(first_command_name));
  cli.AddCommand(&first_command);

  StrictMock<MockCommand> second_command;
  const std::string second_command_name = "second-command";
  const std::vector<const char *> args = { "some-actual-arg" };
  EXPECT_CALL(second_command, Name())
    .WillRepeatedly(ReturnRef(second_command_name));
  EXPECT_CALL(second_command, Run(args))
    .WillOnce(Return(btool::Error::Create("some error")));
  cli.AddCommand(&second_command);

  const char *argv[] = {
    "--some-flag",
    "some-flag-arg",
    "--some-other-flag-arg",
    "some-other-flag-arg",
    "second-command",
    "--some-command-flag",
    "some-command-flag-arg",
    "some-actual-arg",
    "--some-other-command-flag-arg",
    "some-other-command-flag-arg",
  };
  int argc = sizeof(argv)/sizeof(argv[0]);
  EXPECT_EQ(cli.Run(argc, argv), btool::Error::Create("some error"));
}
