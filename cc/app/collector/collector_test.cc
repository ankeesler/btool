#include "collector.h"

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "app/collector/store.h"
#include "err.h"

using ::testing::_;

class MockCollectini : public ::btool::app::collector::Collector::Collectini {
 public:
  MOCK_METHOD1(Collect, void(::btool::app::collector::Store *));
};

TEST(Collector, A) {
  MockCollectini c0;
  EXPECT_CALL(c0, Collect(_));
  MockCollectini c1;
  EXPECT_CALL(c1, Collect(_));
  MockCollectini c2;
  EXPECT_CALL(c2, Collect(_));

  ::btool::app::collector::Store s;
  auto n = s.Put("some/node");

  ::btool::app::collector::Collector c(&s);
  c.AddCollectini(&c0);
  c.AddCollectini(&c1);
  c.AddCollectini(&c2);

  auto err = c.Collect("some/node");
  EXPECT_FALSE(err);
  EXPECT_EQ(n, err.Ret());
}
