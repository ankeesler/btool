#include "app/collector/cc/resolver_factory_delegate.h"

#include <string>
#include <vector>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "app/collector/cc/testing/cc.h"
#include "node/testing/node.h"

using ::testing::Return;
using ::testing::StrictMock;

class ResolverFactoryDelegateTest : public ::testing::Test {
 protected:
  ResolverFactoryDelegateTest() : rfd_(&mrf_) {}

  StrictMock<::btool::node::testing::MockResolver> mr_;
  StrictMock<::btool::app::collector::cc::testing::MockResolverFactory> mrf_;
  ::btool::app::collector::cc::ResolverFactoryDelegate rfd_;
};

TEST_F(ResolverFactoryDelegateTest, CompileC) {
  EXPECT_CALL(mrf_, NewCompileC()).WillOnce(Return(&mr_));

  ::btool::node::Node d1("d1");
  ::btool::node::Node d0("d0");
  d0.dependencies()->push_back(&d1);
  ::btool::node::PropertyStore config;
  ::btool::node::Node n("n");
  n.dependencies()->push_back(&d0);
  EXPECT_EQ(&mr_, rfd_.NewResolver("io.btool.collector.cc.resolvers/compileC",
                                   config, n));
}

TEST_F(ResolverFactoryDelegateTest, CompileCC) {
  EXPECT_CALL(mrf_, NewCompileCC()).WillOnce(Return(&mr_));

  ::btool::node::Node d1("d1");
  ::btool::node::Node d0("d0");
  d0.dependencies()->push_back(&d1);
  ::btool::node::PropertyStore config;
  ::btool::node::Node n("n");
  n.dependencies()->push_back(&d0);
  EXPECT_EQ(&mr_, rfd_.NewResolver("io.btool.collector.cc.resolvers/compileCC",
                                   config, n));
}

TEST_F(ResolverFactoryDelegateTest, LinkC) {
  EXPECT_CALL(mrf_, NewLinkC()).WillOnce(Return(&mr_));

  ::btool::node::Node d1("d1");
  ::btool::node::Node d0("d0");
  d0.dependencies()->push_back(&d1);
  ::btool::node::PropertyStore config;
  ::btool::node::Node n("n");
  n.dependencies()->push_back(&d0);
  EXPECT_EQ(&mr_, rfd_.NewResolver("io.btool.collector.cc.resolvers/linkC",
                                   config, n));
}

TEST_F(ResolverFactoryDelegateTest, LinkCC) {
  EXPECT_CALL(mrf_, NewLinkCC()).WillOnce(Return(&mr_));

  ::btool::node::Node d1("d1");
  ::btool::node::Node d0("d0");
  d0.dependencies()->push_back(&d1);
  ::btool::node::PropertyStore config;
  ::btool::node::Node n("n");
  n.dependencies()->push_back(&d0);
  EXPECT_EQ(&mr_, rfd_.NewResolver("io.btool.collector.cc.resolvers/linkCC",
                                   config, n));
}
