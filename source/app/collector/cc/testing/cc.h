#ifndef BTOOL_APP_COLLECTOR_CC_TESTING_CC_H_
#define BTOOL_APP_COLLECTOR_CC_TESTING_CC_H_

#include <string>
#include <vector>

#include "gmock/gmock.h"

#include "app/collector/cc/resolver_factory.h"
#include "node/node.h"

namespace btool::app::collector::cc::testing {

class MockResolverFactory
    : public ::btool::app::collector::cc::ResolverFactory {
 public:
  MOCK_METHOD0(NewCompileC, class ::btool::node::Node::Resolver *());
  MOCK_METHOD0(NewCompileCC, class ::btool::node::Node::Resolver *());
  MOCK_METHOD0(NewArchive, class ::btool::node::Node::Resolver *());
  MOCK_METHOD0(NewLinkC, class ::btool::node::Node::Resolver *());
  MOCK_METHOD0(NewLinkCC, class ::btool::node::Node::Resolver *());
};

};  // namespace btool::app::collector::cc::testing

#endif  // BTOOL_APP_COLLECTOR_CC_TESTING_CC_H_
