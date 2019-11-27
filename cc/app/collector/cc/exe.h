#ifndef BTOOL_APP_COLLECTOR_CC_EXE_H_
#define BTOOL_APP_COLLECTOR_CC_EXE_H_

#include <string>
#include <vector>

#include "app/collector/base_collectini.h"
#include "app/collector/cc/resolver_factory.h"
#include "app/collector/store.h"
#include "node/node.h"

namespace btool::app::collector::cc {

class Exe : public ::btool::app::collector::BaseCollectini {
 public:
  Exe(::btool::app::collector::cc::ResolverFactory *rf) : rf_(rf) {}

  void OnNotify(::btool::app::collector::Store *s,
                const std::string &name) override;

 private:
  ::btool::app::collector::cc::ResolverFactory *rf_;

  bool CollectObjects(::btool::app::collector::Store *s, ::btool::node::Node *n,
                      const std::string &ext,
                      std::vector<::btool::node::Node *> *objs);
  bool CollectLibraries(::btool::app::collector::Store *s,
                        const ::btool::node::Node &n,
                        std::vector<::btool::node::Node *> *libs);
};

};  // namespace btool::app::collector::cc

#endif  // BTOOL_APP_COLLECTOR_CC_EXE_H_
