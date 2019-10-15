#ifndef BTOOL_APP_COLLECTOR_TESTING_COLLECTOR_H_
#define BTOOL_APP_COLLECTOR_TESTING_COLLECTOR_H_

#include <string>
#include <vector>

#include "gmock/gmock.h"

#include "app/collector/resolver_factory.h"
#include "node/node.h"

namespace btool::app::collector::testing {

class MockResolverFactory : public ::btool::app::collector::ResolverFactory {
 public:
  MOCK_METHOD2(NewCompileC,
               class ::btool::node::Node::Resolver *(std::vector<std::string>,
                                                     std::vector<std::string>));
  MOCK_METHOD2(NewCompileCC,
               class ::btool::node::Node::Resolver *(std::vector<std::string>,
                                                     std::vector<std::string>));
  MOCK_METHOD0(NewArchive, class ::btool::node::Node::Resolver *());
  MOCK_METHOD1(NewLinkC,
               class ::btool::node::Node::Resolver *(std::vector<std::string>));
  MOCK_METHOD1(NewLinkCC,
               class ::btool::node::Node::Resolver *(std::vector<std::string>));

  MOCK_METHOD2(NewDownload,
               class ::btool::node::Node::Resolver *(std::string, std::string));
  MOCK_METHOD0(NewUnzip, class ::btool::node::Node::Resolver *());
};

};  // namespace btool::app::collector::testing

#endif  // BTOOL_APP_COLLECTOR_TESTING_COLLECTOR_H_
