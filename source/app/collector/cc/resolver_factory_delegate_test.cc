#include "app/collector/cc/resolver_factory_delegate.h"

#include <string>
#include <vector>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "app/collector/cc/properties.h"
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
  std::vector<std::string> include_dirs{"include/dir1", "include/dir0"};
  std::vector<std::string> flags{};
  EXPECT_CALL(mrf_, NewCompileC(include_dirs, flags)).WillOnce(Return(&mr_));

  ::btool::node::Node d1("d1");
  ::btool::app::collector::cc::Properties::AddIncludePath(d1.property_store(),
                                                          "include/dir1");
  ::btool::node::Node d0("d0");
  d0.dependencies()->push_back(&d1);
  ::btool::app::collector::cc::Properties::AddIncludePath(d0.property_store(),
                                                          "include/dir0");
  ::btool::node::PropertyStore config;
  ::btool::node::Node n("n");
  n.dependencies()->push_back(&d0);
  EXPECT_EQ(&mr_, rfd_.NewResolver("io.btool.collector.cc.resolvers/compileC",
                                   config, n));
}

TEST_F(ResolverFactoryDelegateTest, CompileCC) {
  std::vector<std::string> include_dirs{"include/dir1", "include/dir0"};
  std::vector<std::string> flags{};
  EXPECT_CALL(mrf_, NewCompileCC(include_dirs, flags)).WillOnce(Return(&mr_));

  ::btool::node::Node d1("d1");
  ::btool::app::collector::cc::Properties::AddIncludePath(d1.property_store(),
                                                          "include/dir1");
  ::btool::node::Node d0("d0");
  d0.dependencies()->push_back(&d1);
  ::btool::app::collector::cc::Properties::AddIncludePath(d0.property_store(),
                                                          "include/dir0");
  ::btool::node::PropertyStore config;
  ::btool::node::Node n("n");
  n.dependencies()->push_back(&d0);
  EXPECT_EQ(&mr_, rfd_.NewResolver("io.btool.collector.cc.resolvers/compileCC",
                                   config, n));
}

TEST_F(ResolverFactoryDelegateTest, LinkC) {
  std::vector<std::string> flags{"-flag1", "-flag0"};
  EXPECT_CALL(mrf_, NewLinkC(flags)).WillOnce(Return(&mr_));

  ::btool::node::Node d1("d1");
  ::btool::app::collector::cc::Properties::AddLinkFlag(d1.property_store(),
                                                       "-flag1");
  ::btool::node::Node d0("d0");
  d0.dependencies()->push_back(&d1);
  ::btool::app::collector::cc::Properties::AddLinkFlag(d0.property_store(),
                                                       "-flag0");
  ::btool::node::PropertyStore config;
  ::btool::node::Node n("n");
  n.dependencies()->push_back(&d0);
  EXPECT_EQ(&mr_, rfd_.NewResolver("io.btool.collector.cc.resolvers/linkC",
                                   config, n));
}

TEST_F(ResolverFactoryDelegateTest, LinkCC) {
  std::vector<std::string> flags{"-flag1", "-flag0"};
  EXPECT_CALL(mrf_, NewLinkCC(flags)).WillOnce(Return(&mr_));

  ::btool::node::Node d1("d1");
  ::btool::app::collector::cc::Properties::AddLinkFlag(d1.property_store(),
                                                       "-flag1");
  ::btool::node::Node d0("d0");
  d0.dependencies()->push_back(&d1);
  ::btool::app::collector::cc::Properties::AddLinkFlag(d0.property_store(),
                                                       "-flag0");
  ::btool::node::PropertyStore config;
  ::btool::node::Node n("n");
  n.dependencies()->push_back(&d0);
  EXPECT_EQ(&mr_, rfd_.NewResolver("io.btool.collector.cc.resolvers/linkCC",
                                   config, n));
}
