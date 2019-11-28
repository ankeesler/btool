#include "app/collector/cc/includes_parser_impl.h"

#include <string>
#include <vector>

#include "gtest/gtest.h"

#include "util/fs/fs.h"

class IncludesParserImplTest : public ::testing::Test {
 protected:
  void SetUp() override { dir_ = ::btool::util::fs::TempDir(); }

  void TearDown() override { ::btool::util::fs::RemoveAll(dir_); }

  std::string dir_;
};

TEST_F(IncludesParserImplTest, Yeah) {
  const std::string file = ::btool::util::fs::Join(dir_, "tuna.h");
  const std::string content =
      "#include \"foo.h\"\n\n"
      "#include <string>\n"
      "#include <cstdio>\n\n"
      "#include \"some/path/to/bar.h\"\n\n"
      "#define IGNORE_THIS\n\n";
  ::btool::util::fs::WriteFile(file, content);

  std::vector<std::string> calls;
  ::btool::app::collector::cc::IncludesParserImpl ipi;
  ipi.ParseIncludes(
      file, [&calls](const std::string &include) { calls.push_back(include); });
  EXPECT_EQ(2UL, calls.size());
  EXPECT_EQ("foo.h", calls[0]);
  EXPECT_EQ("some/path/to/bar.h", calls[1]);
}
