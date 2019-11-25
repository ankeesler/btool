#include "util/unzip_btool.h"

#include "gtest/gtest.h"

#include "util/download.h"
#include "util/fs/fs.h"

TEST(Unzip, Success) {
  auto dir = ::btool::util::fs::TempDir();
  auto file = ::btool::util::fs::Join(dir, "downloaded-stuff");

  ::btool::util::Download("https://github.com/ankeesler/anwork/archive/v9.zip",
                          file);

  ::btool::util::Unzip(file, dir);

  std::size_t dir_count = 0;
  std::size_t file_size_count = 0;
  ::btool::util::fs::Walk(
      dir, [&dir_count, &file_size_count](const std::string &path) {
        if (::btool::util::fs::IsDir(path)) {
          ++dir_count;
        } else {
          auto content = ::btool::util::fs::ReadFile(path);
          file_size_count += content.size();
        }
      });
  EXPECT_EQ(26, dir_count);
  EXPECT_EQ(337994 /* archive file sizes */ + 102844 /* archive size */,
            file_size_count);

  ::btool::util::fs::RemoveAll(dir);
}