#include <unistd.h>
#include <fstream>
#include <sstream>
#include <string>

#include "gtest/gtest.h"
#include "gmock/gmock.h"
#include "fs_impl.h"
#include "error.h"

class FSImplTest : public ::testing::Test {
protected:
  void SetUp() {
    root_ = new char[6] { 'X', 'X', 'X', 'X', 'X', 'X', };
    ASSERT_TRUE((root_ = ::mkdtemp(root_)) != NULL);
  }

  void TearDown() {
    ASSERT_EQ(::rmdir(root_), 0);
    delete root_;
  }

  char *root_;
};

TEST_F(FSImplTest, Success) {
  const std::string root(root_);
  btool::FSImpl fs(root);
  ASSERT_EQ(fs.WriteFile("some_root_dir/some_sub_dir/some_file",
                         "here is a string\nwith a newline"),
            btool::Error::Success());

  std::ifstream file(root + "/some_root_dir/some_sub_dir/some_file");
  ASSERT_TRUE(file);

  std::ostringstream contents;
  contents << file.rdbuf();
  ASSERT_EQ(contents.str(), "here is a string\n with a newline");
}
