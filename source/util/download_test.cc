#include "util/download.h"

#include "gtest/gtest.h"
#include "util/fs/fs.h"

TEST(Download, Success) {
  auto dir = ::btool::util::fs::TempDir();
  auto file = ::btool::util::fs::Join(dir, "downloaded-stuff");

  ::btool::util::Download("https://github.com/ankeesler/anwork/archive/v9.zip",
                          file);

  auto content = ::btool::util::fs::ReadFile(file);
  EXPECT_EQ(102844UL, content.size());

  ::btool::util::fs::RemoveAll(dir);
}
