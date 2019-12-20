#include "app/collector/base_collectini.h"

#include <string>
#include <vector>

#include "gtest/gtest.h"

#include "app/collector/store.h"
#include "app/collector/testing/collector.h"

TEST(BaseCollectini, A) {
  ::btool::app::collector::testing::SpyCollectini a;
  ::btool::app::collector::testing::SpyCollectini b;
  ::btool::app::collector::testing::SpyCollectini c;

  ::btool::app::collector::Store s;
  a.Collect(&s);
  EXPECT_EQ(&s, a.collect_calls_[0]);

  EXPECT_EQ(0UL, a.on_notify_calls_.size());
  EXPECT_EQ(&s, b.on_notify_calls_[0].first);
  EXPECT_EQ("some-other-name", b.on_notify_calls_[0].second);
  EXPECT_EQ(&s, c.on_notify_calls_[0].first);
  EXPECT_EQ("some-other-name", c.on_notify_calls_[0].second);
}
