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

  static const std::vector<std::string> *LinkFlags(
      const ::btool::node::PropertyStore *ps) {
    return ::btool::app::collector::ReadStringsProperty(ps, kLinkFlags);
  }

  static const std::vector<std::string> *Libraries(
      const ::btool::node::PropertyStore *ps) {
    return ::btool::app::collector::ReadStringsProperty(ps, kLibraries);
  }

  static void AddIncludePath(::btool::node::PropertyStore *ps,
                             const std::string &path) {
    ps->Append(kIncludePaths, path);
  }

  static void AddLinkFlag(::btool::node::PropertyStore *ps,
                          const std::string &flag) {
    ps->Append(kLinkFlags, flag);
  }

  static void AddLibrary(::btool::node::PropertyStore *ps,
                         const std::string &library) {
    ps->Append(kLibraries, library);
  }

 private:
  static const char *kIncludePaths;
  static const char *kLinkFlags;
  static const char *kLibraries;
};

};  // namespace btool::app::collector::cc

#endif  // BTOOL_APP_COLLECTOR_CC_PROPERTIES_H_
