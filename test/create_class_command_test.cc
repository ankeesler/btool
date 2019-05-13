#include <string>
#include <vector>

#include "gtest/gtest.h"
#include "gmock/gmock.h"
#include "create_class_command.h"
#include "error.h"
#include "fs.h"

using ::testing::Return;
using ::testing::StrictMock;
using ::testing::_;

class MockFS : public btool::FS {
public:
  MOCK_METHOD2(WriteFile, btool::Error(const std::string&, const std::string&));
};

TEST(CreateClassCommandTest, Success) {
  StrictMock<MockFS> fs;
  EXPECT_CALL(fs, WriteFile(_, _))
    .WillOnce(Return(btool::Error::Success()));
  EXPECT_CALL(fs, WriteFile(_, _))
    .WillOnce(Return(btool::Error::Success()));

  btool::CreateClassCommand command(&fs);

  const std::vector<const char *> args = { "some_rootdir/some_subdir/some_class_name" };
  EXPECT_EQ(command.Run(args), btool::Error::Success());
}
