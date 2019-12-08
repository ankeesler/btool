#include "app/collector/registry/resolver_factory_delegate.h"

#include <string>
#include <vector>

#include "gmock/gmock.h"
#include "gtest/gtest.h"

#include "app/collector/properties.h"
#include "node/node.h"
#include "node/testing/node.h"

using ::testing::Return;
using ::testing::StrictMock;

class MockResolverFactory : public ::btool::app::collector::registry::
                                ResolverFactoryDelegate::ResolverFactory {
 public:
  MOCK_METHOD2(NewDownload,
               ::btool::node::Node::Resolver *(const std::string &,
                                               const std::string &));
  MOCK_METHOD0(NewUnzip, ::btool::node::Node::Resolver *());
  MOCK_METHOD0(NewUntar, ::btool::node::Node::Resolver *());
  MOCK_METHOD3(NewCmd,
               ::btool::node::Node::Resolver *(const std::string &,
                                               const std::vector<std::string> &,
                                               const std::string &));
};

TEST(ResolverFactoryDelegate, Download) {
  StrictMock<MockResolverFactory> mrf;
  StrictMock<::btool::node::testing::MockResolver> mr;
  ::btool::app::collector::registry::ResolverFactoryDelegate rfd(&mrf);
  EXPECT_CALL(mrf, NewDownload("some-url", "some-sha256"))
      .WillOnce(Return(&mr));

  ::btool::node::PropertyStore config;
  config.Write("url", "some-url");
  config.Write("sha256", "some-sha256");
  ::btool::node::Node n("n");
  EXPECT_EQ(&mr, rfd.NewResolver("download", config, n));
}

TEST(ResolverFactoryDelegate, Unzip) {
  StrictMock<MockResolverFactory> mrf;
  StrictMock<::btool::node::testing::MockResolver> mr;
  ::btool::app::collector::registry::ResolverFactoryDelegate rfd(&mrf);
  EXPECT_CALL(mrf, NewUnzip()).WillOnce(Return(&mr));

  ::btool::node::PropertyStore config;
  ::btool::node::Node n("n");
  EXPECT_EQ(&mr, rfd.NewResolver("unzip", config, n));
}

TEST(ResolverFactoryDelegate, Cmd) {
  StrictMock<MockResolverFactory> mrf;
  StrictMock<::btool::node::testing::MockResolver> mr;
  ::btool::app::collector::registry::ResolverFactoryDelegate rfd(&mrf);
  std::vector<std::string> args{"some-arg-0", "some-arg-1"};
  EXPECT_CALL(mrf, NewCmd("some-path", args, "some-root/some-dir"))
      .WillOnce(Return(&mr));

  ::btool::node::PropertyStore config;
  config.Write("path", "some-path");
  config.Append("args", "some-arg-0");
  config.Append("args", "some-arg-1");
  config.Write("dir", "some-dir");
  ::btool::node::Node n("n");
  ::btool::app::collector::Properties::SetRoot(n.property_store(), "some-root");
  EXPECT_EQ(&mr, rfd.NewResolver("cmd", config, n));
}
