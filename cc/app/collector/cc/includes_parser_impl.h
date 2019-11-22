#ifndef BTOOL_APP_COLLECTOR_CC_INCLUDESPARSERIMPL_H_
#define BTOOL_APP_COLLECTOR_CC_INCLUDESPARSERIMPL_H_

#include <functional>
#include <string>

#include "app/collector/cc/inc.h"

namespace btool::app::collector::cc {

class IncludesParserImpl : public Inc::IncludesParser {
 public:
  void ParseIncludes(
      const std::string &path,
      std::function<void(const std::string &)> callback) override;
};

};  // namespace btool::app::collector::cc

#endif  // BTOOL_APP_COLLECTOR_CC_INCLUDESPARSERIMPL_H_
