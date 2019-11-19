#ifndef BTOOL_APP_COLLECTOR_CC_INC_H_
#define BTOOL_APP_COLLECTOR_CC_INC_H_

#include <functional>
#include <map>
#include <string>
#include <vector>

#include "app/collector/base_collectini.h"
#include "app/collector/store.h"
#include "err.h"

namespace btool::app::collector::cc {

class Inc : public ::btool::app::collector::BaseCollectini {
 public:
  class IncludesParser {
   public:
    virtual ~IncludesParser() {}
    virtual ::btool::VoidErr ParseIncludes(
        const std::string &path,
        std::function<void(const std::string &)> callback) = 0;
  };

  Inc(IncludesParser *ip) : ip_(ip) {}

  void OnNotify(::btool::app::collector::Store *s,
                const std::string &name) override;

 private:
  IncludesParser *ip_;
};

};  // namespace btool::app::collector::cc

#endif  // BTOOL_APP_COLLECTOR_CC_INC_H_
