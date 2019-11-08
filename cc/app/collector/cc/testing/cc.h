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
  MOCK_METHOD2(NewCompileC, class ::btool::node::Node::Resolver *(
                                const std::vector<std::string> &,
                                const std::vector<std::string> &));
  MOCK_METHOD2(NewCompileCC, class ::btool::node::Node::Resolver *(
                                 const std::vector<std::string> &,
                                 const std::vector<std::string> &));
  MOCK_METHOD0(NewArchive, class ::btool::node::Node::Resolver *());
  MOCK_METHOD1(NewLinkC, class ::btool::node::Node::Resolver *(
                             const std::vector<std::string> &));
  MOCK_METHOD1(NewLinkCC, class ::btool::node::Node::Resolver *(
                              const std::vector<std::string> &));
};

};  // namespace btool::app::collector::cc::testing

#endif  // BTOOL_APP_COLLECTOR_CC_TESTING_CC_H_
