#ifndef BTOOL_APP_COLLECTOR_TESTING_COLLECTOR_H_
#define BTOOL_APP_COLLECTOR_TESTING_COLLECTOR_H_

#include <string>
#include <vector>

#include "gmock/gmock.h"

#include "app/collector/base_collectini.h"
#include "app/collector/resolver_factory_delegate.h"
#include "node/node.h"

namespace btool::app::collector::testing {

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

class MockResolverFactoryDelegate
    : public ::btool::app::collector::ResolverFactoryDelegate {
 public:
  MOCK_METHOD4(NewResolver,
               ::btool::node::Node::Resolver *(
                   const std::string &, const ::btool::node::PropertyStore &,
                   const std::string &, const ::btool::node::Node &));
};

};  // namespace btool::app::collector::testing

#endif  // BTOOL_APP_COLLECTOR_TESTING_COLLECTOR_H_
