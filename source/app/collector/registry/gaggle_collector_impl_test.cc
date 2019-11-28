#include "app/collector/registry/gaggle_collector_impl.h"

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "app/collector/registry/registry.h"
#include "app/collector/store.h"
#include "app/collector/testing/collector.h"
#include "node/node.h"
#include "node/testing/node.h"

using ::testing::_;
using ::testing::ElementsAre;
using ::testing::InSequence;
using ::testing::Ref;
using ::testing::Return;
using ::testing::StrictMock;

// class MockResolverFactory : public ::btool::app::collector::registry::
//                                GaggleCollectorImpl::ResolverFactory {
// public:
//  MOCK_METHOD1(New, ::btool::node::Node::Resolver *(
//                        const ::btool::app::collector::registry::Resolver &));
//};

class GaggleCollectorImplTest : public ::testing::Test {
 protected:
  void SetUp() override {
    gci_.AddResolverFactoryDelegate(&mrfd0_);
    gci_.AddResolverFactoryDelegate(&mrfd1_);

    n0_.labels.Write("bool-property", true);
    n0_.labels.Append("strings-property", "some-string");
    g_.nodes = {n0_, n1_};
  }

  // StrictMock<MockResolverFactory> mrf_;
  StrictMock<::btool::app::collector::testing::MockResolverFactoryDelegate>
      mrfd0_;
  StrictMock<::btool::app::collector::testing::MockResolverFactoryDelegate>
      mrfd1_;
  ::btool::app::collector::registry::GaggleCollectorImpl gci_;
  ::btool::app::collector::registry::Resolver r0_{.name = "r0"};
  ::btool::app::collector::registry::Resolver r1_{.name = "r1"};
  ::btool::app::collector::registry::Node n0_{.name = "n0", .resolver = r0_};
  ::btool::app::collector::registry::Node n1_{
      .name = "n1", .dependencies = {"n0"}, .resolver = r1_};
  ::btool::app::collector::registry::Gaggle g_;
  ::btool::app::collector::Store s_;
};

TEST_F(GaggleCollectorImplTest, Success) {
  InSequence s;

  StrictMock<::btool::node::testing::MockResolver> mr0;
  StrictMock<::btool::node::testing::MockResolver> mr1;
  EXPECT_CALL(mrfd0_, NewResolver(r0_.name, r0_.config, _))
      .WillOnce(Return(&mr0));
  EXPECT_CALL(mrfd0_, NewResolver(r1_.name, r1_.config, _))
      .WillOnce(Return(nullptr));
  EXPECT_CALL(mrfd1_, NewResolver(r1_.name, r1_.config, _))
      .WillOnce(Return(&mr1));

  gci_.Collect(&s_, &g_, "some-root");

  auto n0 = s_.Get("some-root/n0");
  ASSERT_TRUE(n0 != nullptr);
  EXPECT_TRUE(n0->dependencies()->empty());
  EXPECT_EQ(&mr0, n0->resolver());

  const bool *b;
  n0->property_store()->Read("bool-property", &b);
  ASSERT_TRUE(b != nullptr);
  EXPECT_TRUE(*b);

  const std::vector<std::string> *ss;
  n0->property_store()->Read("strings-property", &ss);
  ASSERT_TRUE(ss != nullptr);
  EXPECT_THAT(*ss, ElementsAre("some-string"));

  auto n1 = s_.Get("some-root/n1");
  EXPECT_TRUE(n1 != nullptr);
  EXPECT_THAT(*n1->dependencies(), ElementsAre(n0));
  EXPECT_EQ(&mr1, n1->resolver());
}
