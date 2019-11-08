#ifndef BTOOL_APP_COLLECTOR_TESTING_COLLECTOR_H_
#define BTOOL_APP_COLLECTOR_TESTING_COLLECTOR_H_

#include <string>
#include <vector>

#include "gmock/gmock.h"

#include "app/collector/base_collectini.h"
#include "app/collector/resolver_factory.h"
#include "node/node.h"

namespace btool::app::collector::testing {

class MockResolverFactory : public ::btool::app::collector::ResolverFactory {
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

  MOCK_METHOD2(NewDownload,
               class ::btool::node::Node::Resolver *(const std::string &,
                                                     const std::string &));
  MOCK_METHOD0(NewUnzip, class ::btool::node::Node::Resolver *());
};

class SpyCollectini : public ::btool::app::collector::BaseCollectini {
 public:
  std::vector<::btool::app::collector::Store *> collect_calls_;
  std::vector<std::pair<::btool::app::collector::Store *, std::string>>
      on_notify_calls_;

  void Collect(::btool::app::collector::Store *s) override {
    collect_calls_.push_back(s);
    Notify(s, "some-other-name");
  }

 protected:
  void OnNotify(::btool::app::collector::Store *s,
                const std::string &name) override {
    on_notify_calls_.push_back({s, name});
  }
};

};  // namespace btool::app::collector::testing

#endif  // BTOOL_APP_COLLECTOR_TESTING_COLLECTOR_H_
