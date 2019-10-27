#ifndef BTOOL_APP_COLLECTOR_CC_EXE_H_
#define BTOOL_APP_COLLECTOR_CC_EXE_H_

#include <string>

#include "app/collector/listener_collectini.h"
#include "app/collector/resolver_factory.h"
#include "app/collector/store.h"

namespace btool::app::collector::cc {

class Exe : public ::btool::app::collector::ListenerCollectini {
 public:
  Exe(::btool::app::collector::ResolverFactory *rf) : rf_(rf) {}

  void OnSet(::btool::app::collector::Store *s,
             const std::string &name) override;

 private:
  ::btool::app::collector::ResolverFactory *rf_;
};

};  // namespace btool::app::collector::cc

#endif  // BTOOL_APP_COLLECTOR_CC_EXE_H_
