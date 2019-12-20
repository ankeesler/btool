#include "app/collector/cc/includes_parser_impl.h"

#include <string>
#include <utility>
#include <vector>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "util/fs/fs.h"

using ::testing::ElementsAre;

class IncludesParserImplTest : public ::testing::Test {
 protected:
  void SetUp() override { dir_ = ::btool::util::fs::TempDir(); }

  void TearDown() override { ::btool::util::fs::RemoveAll(dir_); }

  std::string dir_;
};

TEST_F(IncludesParserImplTest, Success) {
  const std::string file = ::btool::util::fs::Join(dir_, "tuna.h");
  const std::string content =
      "#include \"foo.h\"\n\n"
      "#include <string>\n"
      "#include <cstdio>\n\n"
      "#include \"some/path/to/bar.h\"\n\n"
      "#define IGNORE_THIS\n\n"
      "#include \"comment.h\" // here is a comment\n"
      "#include \"one/more.h\"\n";
  ::btool::util::fs::WriteFile(file, content);

  std::vector<std::pair<std::string, bool>> calls;
  ::btool::app::collector::cc::IncludesParserImpl ipi;
  ipi.ParseIncludes(file, [&calls](const std::string &include, bool system) {
    calls.push_back(std::pair<std::string, bool>(include, system));
  });
  EXPECT_THAT(
      calls, ElementsAre(
                 std::make_pair<std::string, bool>("foo.h", false),
                 std::make_pair<std::string, bool>("string", true),
                 std::make_pair<std::string, bool>("cstdio", true),
                 std::make_pair<std::string, bool>("some/path/to/bar.h", false),
                 std::make_pair<std::string, bool>("comment.h", false),
                 std::make_pair<std::string, bool>("one/more.h", false)));
}

TEST_F(IncludesParserImplTest, EmptyString) {
  const std::string file = ::btool::util::fs::Join(dir_, "tuna.h");
  const std::string content =
      "#include <iostream>\n"
      "#include <ostream>\n"
      "#include <sstream>\n"
      "#include <string>\n"
      "#include <vector>\n\n"
      "const std::string empty = \"\";\n";
  ::btool::util::fs::WriteFile(file, content);

  std::vector<std::string> calls;
  ::btool::app::collector::cc::IncludesParserImpl ipi;
  ipi.ParseIncludes(file, [&calls](const std::string &include, bool system) {
    calls.push_back(include);
  });
  EXPECT_THAT(
      calls, ElementsAre("iostream", "ostream", "sstream", "string", "vector"));
}
