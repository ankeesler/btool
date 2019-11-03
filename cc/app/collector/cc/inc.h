#ifndef BTOOL_APP_COLLECTOR_CC_INC_H_
#define BTOOL_APP_COLLECTOR_CC_INC_H_

#include <functional>
#include <string>

#include "app/collector/listener_collectini.h"
#include "app/collector/store.h"
#include "core/err.h"

namespace btool::app::collector::cc {

class Inc : public ::btool::app::collector::ListenerCollectini {
 public:
  class IncludesParser {
   public:
    virtual ~IncludesParser() {}
    virtual ::btool::core::VoidErr ParseIncludes(
        const std::string &path,
        std::function<void(const std::string &)> callback) = 0;
  };

  Inc(IncludesParser *ip) : ip_(ip) {}

  void OnSet(::btool::app::collector::Store *s,
             const std::string &name) override;

 private:
  IncludesParser *ip_;
};

};  // namespace btool::app::collector::cc

#endif  // BTOOL_APP_COLLECTOR_CC_INC_H_
