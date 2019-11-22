#include "app.h"

#include <string>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "err.h"
#include "node/node.h"

using ::testing::_;
using ::testing::DoAll;
using ::testing::InSequence;
using ::testing::Ref;
using ::testing::Return;
using ::testing::SetArgPointee;
using ::testing::StrictMock;

class MockCollector : public ::btool::app::App::Collector {
 public:
  MOCK_METHOD3(Collect, bool(const std::string &, ::btool::node::Node **,
                             std::string *));
};

class MockCleaner : public ::btool::app::App::Cleaner {
 public:
  MOCK_METHOD2(Clean, bool(const ::btool::node::Node &, std::string *));
};

class MockLister : public ::btool::app::App::Lister {
 public:
  MOCK_METHOD2(List, bool(const ::btool::node::Node &, std::string *));
};

class MockBuilder : public ::btool::app::App::Builder {
 public:
  MOCK_METHOD2(Build, bool(const ::btool::node::Node &, std::string *));
};

class MockRunner : public ::btool::app::App::Runner {
 public:
  MOCK_METHOD2(Run, bool(const ::btool::node::Node &, std::string *));
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
  auto n = new ::btool::node::Node("a");
  EXPECT_CALL(mcollector_, Collect(_, _, _))
      .WillOnce(DoAll(SetArgPointee<1>(n), Return(true)));
  EXPECT_CALL(mbuilder_, Build(Ref(*n), _)).WillOnce(Return(true));

  std::string err;
  EXPECT_TRUE(a_.Run("", false, false, false, &err)) << "error: " << err;
}

TEST_F(AppTest, Clean) {
  InSequence s;
  auto n = new ::btool::node::Node("a");
  EXPECT_CALL(mcollector_, Collect(_, _, _))
      .WillOnce(DoAll(SetArgPointee<1>(n), Return(true)));
  EXPECT_CALL(mcleaner_, Clean(Ref(*n), _)).WillOnce(Return(true));

  std::string err;
  EXPECT_TRUE(a_.Run("", true, false, false, &err)) << "error: " << err;
}

TEST_F(AppTest, List) {
  InSequence s;
  auto n = new ::btool::node::Node("a");
  EXPECT_CALL(mcollector_, Collect(_, _, _))
      .WillOnce(DoAll(SetArgPointee<1>(n), Return(true)));
  EXPECT_CALL(mlister_, List(Ref(*n), _)).WillOnce(Return(true));

  std::string err;
  EXPECT_TRUE(a_.Run("", false, true, false, &err)) << "error: " << err;
}

TEST_F(AppTest, Run) {
  InSequence s;
  auto n = new ::btool::node::Node("a");
  EXPECT_CALL(mcollector_, Collect(_, _, _))
      .WillOnce(DoAll(SetArgPointee<1>(n), Return(true)));
  EXPECT_CALL(mbuilder_, Build(Ref(*n), _)).WillOnce(Return(true));
  EXPECT_CALL(mrunner_, Run(Ref(*n), _)).WillOnce(Return(true));

  std::string err;
  EXPECT_TRUE(a_.Run("", false, false, true, &err)) << "error: " << err;
}
