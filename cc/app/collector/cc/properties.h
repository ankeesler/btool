#ifndef BTOOL_APP_COLLECTOR_CC_PROPERTIES_H_
#define BTOOL_APP_COLLECTOR_CC_PROPERTIES_H_

#include <string>
#include <vector>

#include "app/collector/properties.h"
#include "node/property_store.h"

namespace btool::app::collector::cc {

class Properties {
 public:
  static const std::vector<std::string> *IncludePaths(
      const ::btool::node::PropertyStore *ps) {
    return ::btool::app::collector::ReadStringsProperty(ps, kIncludePaths);
  }

  static void AddIncludePath(::btool::node::PropertyStore *ps,
                             const std::string &path) {
    ps->Append(kIncludePaths, path);
  }

 private:
  static const std::string kIncludePaths;
};

};  // namespace btool::app::collector::cc

#endif  // BTOOL_APP_COLLECTOR_CC_PROPERTIES_H_
