#include "collector.h"

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "core/err.h"
#include "node/store.h"

using ::testing::_;
using ::testing::Return;

class MockCollectini : public ::btool::app::collector::Collector::Collectini {
 public:
  MOCK_METHOD1(Collect, ::btool::core::VoidErr(::btool::node::Store *));
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
  std::vector<::btool::app::collector::Collector::Collectini *> cs{&c0, &c1,
                                                                   &c2};
  ::btool::app::collector::Collector c(&cs);
  EXPECT_FALSE(c.Collect());
}
