#include "collector.h"

#include <string>
#include <vector>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "app/collector/store.h"

using ::testing::_;
using ::testing::HasSubstr;
using ::testing::Return;

class MockCollectini : public ::btool::app::collector::Collector::Collectini {
 public:
  MOCK_METHOD1(Collect, void(::btool::app::collector::Store *));
  MOCK_METHOD0(Errors, std::vector<std::string>());
};

TEST(Collector, Pass) {
  MockCollectini c0;
  EXPECT_CALL(c0, Collect(_));
  EXPECT_CALL(c0, Errors());
  MockCollectini c1;
  EXPECT_CALL(c1, Collect(_));
  EXPECT_CALL(c1, Errors());
  MockCollectini c2;
  EXPECT_CALL(c2, Collect(_));
  EXPECT_CALL(c2, Errors());

  ::btool::app::collector::Store s;
  auto n = s.Put("some/node");

  ::btool::app::collector::Collector c(&s);
  c.AddCollectini(&c0);
  c.AddCollectini(&c1);
  c.AddCollectini(&c2);

  ::btool::node::Node *ret_n;
  std::string ret_err;
  EXPECT_TRUE(c.Collect("some/node", &ret_n, &ret_err)) << "error: " << ret_err;
  EXPECT_EQ(n, ret_n);
}

TEST(Collector, Fail) {
  std::vector<std::string> errors{"some error"};

  MockCollectini c0;
  EXPECT_CALL(c0, Collect(_));
  EXPECT_CALL(c0, Errors());
  MockCollectini c1;
  EXPECT_CALL(c1, Collect(_));
  EXPECT_CALL(c1, Errors()).WillOnce(Return(errors));
  MockCollectini c2;

  ::btool::app::collector::Store s;
  s.Put("some/node");

  ::btool::app::collector::Collector c(&s);
  c.AddCollectini(&c0);
  c.AddCollectini(&c1);
  c.AddCollectini(&c2);

  ::btool::node::Node *ret_n;
  std::string ret_err;
  EXPECT_FALSE(c.Collect("some/node", &ret_n, &ret_err))
      << "error: " << ret_err;
  EXPECT_THAT(ret_err, HasSubstr("collect errors:\nsome error"));
}
