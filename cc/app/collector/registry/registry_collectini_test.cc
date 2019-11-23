#include "app/collector/registry/registry_collectini.h"

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "app/collector/registry/registry.h"
#include "app/collector/registry/testing/registry.h"
#include "app/collector/store.h"

using ::testing::_;
using ::testing::DoAll;
using ::testing::InSequence;
using ::testing::SetArgPointee;
using ::testing::StrictMock;

class MockGaggleCollector : public ::btool::app::collector::registry::
                                RegistryCollectini::GaggleCollector {
 public:
  MOCK_METHOD3(Collect,
               void(::btool::app::collector::Store *,
                    ::btool::app::collector::registry::Gaggle *, std::string));
};

TEST(RegistryCollectini, Basic) {
  InSequence seq;

  ::btool::app::collector::registry::IndexFile if0{.path = "some-path-0",
                                                   .sha256 = "sha0"};
  ::btool::app::collector::registry::IndexFile if1{.path = "some-path-1",
                                                   .sha256 = "sha1"};
  ::btool::app::collector::registry::Index i{.files = {if0, if1}};
  ::btool::app::collector::registry::Gaggle g0{};
  ::btool::app::collector::registry::Gaggle g1{};
  std::string cache = "some-cache";
  ::btool::app::collector::Store s;

  StrictMock<::btool::app::collector::registry::testing::MockRegistry> mr;
  StrictMock<MockGaggleCollector> mgc;
  EXPECT_CALL(mr, GetIndex(_)).WillOnce(SetArgPointee<0>(i));
  EXPECT_CALL(mr, GetGaggle("some-path-0", _)).WillOnce(SetArgPointee<1>(g0));
  EXPECT_CALL(mgc, Collect(&s, _, "some-cache/sha0"));
  EXPECT_CALL(mr, GetGaggle("some-path-1", _)).WillOnce(SetArgPointee<1>(g1));
  EXPECT_CALL(mgc, Collect(&s, _, "some-cache/sha1"));

  ::btool::app::collector::registry::RegistryCollectini rc(&mr, cache, &mgc);
  rc.Collect(&s);
}
