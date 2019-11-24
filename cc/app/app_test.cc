#include "app.h"

#include <memory>
#include <string>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "err.h"
#include "node/node.h"

using ::testing::_;
using ::testing::InSequence;
using ::testing::Return;
using ::testing::StrictMock;

class MockCollector : public ::btool::app::App::Collector {
 public:
  MOCK_METHOD1(Collect, ::btool::node::Node *(const std::string &));
};

class MockCleaner : public ::btool::app::App::Cleaner {
 public:
  MOCK_METHOD1(Clean, void(const ::btool::node::Node &));
};

class MockLister : public ::btool::app::App::Lister {
 public:
  MOCK_METHOD1(List, void(const ::btool::node::Node &));
};

class MockBuilder : public ::btool::app::App::Builder {
 public:
  MOCK_METHOD1(Build, void(const ::btool::node::Node &));
};

class MockRunner : public ::btool::app::App::Runner {
 public:
  MOCK_METHOD1(Run, void(const ::btool::node::Node &));
};

class AppTest : public ::testing::Test {
 protected:
  AppTest() : a_(&mcollector_, &mcleaner_, &mlister_, &mbuilder_, &mrunner_) {}

  ::btool::app::App a_;
  StrictMock<MockCollector> mcollector_;
  StrictMock<MockCleaner> mcleaner_;
  StrictMock<MockLister> mlister_;
  StrictMock<MockBuilder> mbuilder_;
  StrictMock<MockRunner> mrunner_;
};

TEST_F(AppTest, Build) {
  InSequence s;
  auto n = std::make_unique<::btool::node::Node>("a");
  EXPECT_CALL(mcollector_, Collect(_)).WillOnce(Return(n.get()));
  EXPECT_CALL(mbuilder_, Build(_));

  a_.Run("", false, false, false);
}

TEST_F(AppTest, Clean) {
  InSequence s;
  auto n = std::make_unique<::btool::node::Node>("a");
  EXPECT_CALL(mcollector_, Collect(_)).WillOnce(Return(n.get()));
  EXPECT_CALL(mcleaner_, Clean(_));

  a_.Run("", true, false, false);
}

TEST_F(AppTest, List) {
  InSequence s;
  auto n = std::make_unique<::btool::node::Node>("a");
  EXPECT_CALL(mcollector_, Collect(_)).WillOnce(Return(n.get()));
  EXPECT_CALL(mlister_, List(_));

  a_.Run("", false, true, false);
}

TEST_F(AppTest, Run) {
  InSequence s;
  auto n = std::make_unique<::btool::node::Node>("a");
  EXPECT_CALL(mcollector_, Collect(_)).WillOnce(Return(n.get()));
  EXPECT_CALL(mbuilder_, Build(_));
  EXPECT_CALL(mrunner_, Run(_));

  a_.Run("", false, false, true);
}
