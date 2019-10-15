#include "collector.h"

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "app/collector/store.h"
#include "core/err.h"

using ::testing::_;
using ::testing::Return;

class MockCollectini : public ::btool::app::collector::Collector::Collectini {
 public:
  MOCK_METHOD1(Collect,
               ::btool::core::VoidErr(::btool::app::collector::Store *));
};

TEST(Collector, A) {
  MockCollectini c0;
  EXPECT_CALL(c0, Collect(_))
      .WillOnce(Return(::btool::core::VoidErr::Success()));
  MockCollectini c1;
  EXPECT_CALL(c1, Collect(_))
      .WillOnce(Return(::btool::core::VoidErr::Success()));
  MockCollectini c2;
  EXPECT_CALL(c2, Collect(_))
      .WillOnce(Return(::btool::core::VoidErr::Success()));

  ::btool::app::collector::Collector c;
  c.AddCollectini(&c0);
  c.AddCollectini(&c1);
  c.AddCollectini(&c2);

  EXPECT_FALSE(c.Collect());
}
