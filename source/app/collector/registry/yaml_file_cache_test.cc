#include "app/collector/registry/yaml_file_cache.h"

#include <string>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "app/collector/registry/testing/registry.h"

using ::testing::_;
using ::testing::InSequence;
using ::testing::StrictMock;

TEST(YamlFileCache, Success) {
  std::string dir = ::btool::util::fs::TempDir();

  InSequence seq;
  StrictMock<
      ::btool::app::collector::registry::testing::MockSerializer<std::string>>
      ms;
  EXPECT_CALL(ms, Marshal(_, "some-value"));
  EXPECT_CALL(ms, Unmarshal(_, _));

  ::btool::app::collector::registry::YamlFileCache yfc(&ms, dir);

  std::string s = "";
  EXPECT_FALSE(yfc.Get("some-name", &s));
  yfc.Set("some-name", "some-value");
  EXPECT_TRUE(yfc.Get("some-name", &s));
  EXPECT_FALSE(yfc.Get("some-other-name", &s));

  ::btool::util::fs::RemoveAll(dir);
}
