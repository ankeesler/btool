#ifndef BTOOL_APP_COLLECTOR_CC_EXE_H_
#define BTOOL_APP_COLLECTOR_CC_EXE_H_

#include <string>

#include "app/collector/base_collectini.h"
#include "app/collector/cc/resolver_factory.h"
#include "app/collector/store.h"

namespace btool::app::collector::cc {

class Exe : public ::btool::app::collector::BaseCollectini {
 public:
  Exe(::btool::app::collector::cc::ResolverFactory *rf) : rf_(rf) {}

  void Collect(::btool::app::collector::Store *s) override {}
  void OnNotify(::btool::app::collector::Store *s,
                const std::string &name) override;

 private:
  ::btool::app::collector::cc::ResolverFactory *rf_;
};

};  // namespace btool::app::collector::cc

#endif  // BTOOL_APP_COLLECTOR_CC_EXE_H_
