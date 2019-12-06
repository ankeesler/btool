#include "app/collector/registry/yaml_file_cache.h"

#include <chrono>
#include <string>
#include <thread>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "app/collector/registry/testing/registry.h"

using ::testing::_;
using ::testing::InSequence;
using ::testing::StrictMock;

TEST(YamlFileCache, LongTimeout) {
  std::string dir = ::btool::util::fs::TempDir();

  InSequence seq;
  StrictMock<
      ::btool::app::collector::registry::testing::MockSerializer<std::string>>
      ms;
  EXPECT_CALL(ms, Marshal(_, "some-value"));
  EXPECT_CALL(ms, Unmarshal(_, _));

  ::btool::app::collector::registry::YamlFileCache yfc(
      &ms, dir, std::chrono::seconds(10));

  std::string s = "";
  EXPECT_FALSE(yfc.Get("some-name", &s));
  yfc.Set("some-name", "some-value");
  EXPECT_TRUE(yfc.Get("some-name", &s));
  EXPECT_FALSE(yfc.Get("some-other-name", &s));

  ::btool::util::fs::RemoveAll(dir);
}

TEST(YamlFileCache, ShortTimeout) {
  std::string dir = ::btool::util::fs::TempDir();

  InSequence seq;
  StrictMock<
      ::btool::app::collector::registry::testing::MockSerializer<std::string>>
      ms;
  EXPECT_CALL(ms, Marshal(_, "some-value"));
  EXPECT_CALL(ms, Unmarshal(_, _));

  std::chrono::seconds timeout(1);
  ::btool::app::collector::registry::YamlFileCache yfc(&ms, dir, timeout);

  std::string s = "";
  yfc.Set("some-name", "some-value");
  EXPECT_TRUE(yfc.Get("some-name", &s));
  std::this_thread::sleep_for(timeout);
  EXPECT_FALSE(yfc.Get("some-name", &s));

  ::btool::util::fs::RemoveAll(dir);
}

TEST(YamlFileCache, NoTimeout) {
  std::string dir = ::btool::util::fs::TempDir();

  InSequence seq;
  StrictMock<
      ::btool::app::collector::registry::testing::MockSerializer<std::string>>
      ms;
  EXPECT_CALL(ms, Marshal(_, "some-value"));

  ::btool::app::collector::registry::YamlFileCache yfc(&ms, dir,
                                                       std::chrono::seconds(0));

  std::string s = "";
  yfc.Set("some-name", "some-value");
  EXPECT_FALSE(yfc.Get("some-name", &s));

  ::btool::util::fs::RemoveAll(dir);
}
