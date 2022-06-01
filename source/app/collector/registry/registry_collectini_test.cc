#include "app/collector/registry/registry_collectini.h"

#include "app/collector/registry/registry.h"
#include "app/collector/registry/testing/registry.h"
#include "app/collector/store.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "util/testing/util.h"

using ::testing::_;
using ::testing::DoAll;
using ::testing::InSequence;
using ::testing::Return;
using ::testing::SetArgPointee;
using ::testing::StrictMock;

class MockGaggleCollector : public ::btool::app::collector::registry::
                                RegistryCollectini::GaggleCollector {
 public:
  MOCK_METHOD3(Collect,
               void(::btool::app::collector::Store *,
                    ::btool::app::collector::registry::Gaggle *, std::string));
};

class RegistryCollectiniTest : public ::testing::Test {
 protected:
  RegistryCollectiniTest() : rc_(&mr_, cache_, &mc_i_, &mc_g_, &mgc_) {}

  void TearDown() override { rc_.Collect(&s_); }

  ::btool::app::collector::registry::IndexFile if0_{.path = "some-path-0",
                                                    .sha256 = "sha0"};
  ::btool::app::collector::registry::IndexFile if1_{.path = "some-path-1",
                                                    .sha256 = "sha1"};
  ::btool::app::collector::registry::Index i_{.files = {if0_, if1_}};
  ::btool::app::collector::registry::Gaggle g0_{.nodes = {}};
  ::btool::app::collector::registry::Gaggle g1_{.nodes = {}};
  std::string cache_ = "some-cache";
  ::btool::app::collector::Store s_;

  StrictMock<::btool::app::collector::registry::testing::MockRegistry> mr_;
  StrictMock<::btool::util::testing::MockCache<
      ::btool::app::collector::registry::Index>>
      mc_i_;
  StrictMock<::btool::util::testing::MockCache<
      ::btool::app::collector::registry::Gaggle>>
      mc_g_;
  StrictMock<MockGaggleCollector> mgc_;

  ::btool::app::collector::registry::RegistryCollectini rc_;
};

TEST_F(RegistryCollectiniTest, NoCache) {
  InSequence seq;

  EXPECT_CALL(mr_, GetName()).WillOnce(Return("some-index"));
  EXPECT_CALL(mc_i_, Get("3357cd63016f9e9a2155a9cb135a06a676d6f3e6d361f6ec9226b"
                         "8645f695f90/index",
                         _))
      .WillOnce(Return(false));
  EXPECT_CALL(mr_, GetIndex(_)).WillOnce(SetArgPointee<0>(i_));
  EXPECT_CALL(mc_i_, Set("3357cd63016f9e9a2155a9cb135a06a676d6f3e6d361f6ec9226b"
                         "8645f695f90/index",
                         i_));

  EXPECT_CALL(mc_g_, Get("3357cd63016f9e9a2155a9cb135a06a676d6f3e6d361f6ec9226b"
                         "8645f695f90/sha0",
                         _))
      .WillOnce(Return(false));
  EXPECT_CALL(mr_, GetGaggle("some-path-0", _)).WillOnce(SetArgPointee<1>(g0_));
  EXPECT_CALL(mc_g_, Set("3357cd63016f9e9a2155a9cb135a06a676d6f3e6d361f6ec9226b"
                         "8645f695f90/sha0",
                         g0_));
  EXPECT_CALL(mgc_, Collect(&s_, _, "some-cache/sha0"));

  EXPECT_CALL(mc_g_, Get("3357cd63016f9e9a2155a9cb135a06a676d6f3e6d361f6ec9226b"
                         "8645f695f90/sha1",
                         _))
      .WillOnce(Return(false));
  EXPECT_CALL(mr_, GetGaggle("some-path-1", _)).WillOnce(SetArgPointee<1>(g1_));
  EXPECT_CALL(mc_g_, Set("3357cd63016f9e9a2155a9cb135a06a676d6f3e6d361f6ec9226b"
                         "8645f695f90/sha1",
                         g1_));
  EXPECT_CALL(mgc_, Collect(&s_, _, "some-cache/sha1"));
}

TEST_F(RegistryCollectiniTest, Cache) {
  InSequence seq;

  EXPECT_CALL(mr_, GetName()).WillOnce(Return("some-index"));
  EXPECT_CALL(mc_i_, Get("3357cd63016f9e9a2155a9cb135a06a676d6f3e6d361f6ec9226b"
                         "8645f695f90/index",
                         _))
      .WillOnce(DoAll(SetArgPointee<1>(i_), Return(true)));

  EXPECT_CALL(mc_g_, Get("3357cd63016f9e9a2155a9cb135a06a676d6f3e6d361f6ec9226b"
                         "8645f695f90/sha0",
                         _))
      .WillOnce(DoAll(SetArgPointee<1>(g0_), Return(true)));
  EXPECT_CALL(mgc_, Collect(&s_, _, "some-cache/sha0"));

  EXPECT_CALL(mc_g_, Get("3357cd63016f9e9a2155a9cb135a06a676d6f3e6d361f6ec9226b"
                         "8645f695f90/sha1",
                         _))
      .WillOnce(DoAll(SetArgPointee<1>(g1_), Return(true)));
  EXPECT_CALL(mgc_, Collect(&s_, _, "some-cache/sha1"));
}
