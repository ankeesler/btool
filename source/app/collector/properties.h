#ifndef BTOOL_APP_COLLECTOR_PROPERTIES_H_
#define BTOOL_APP_COLLECTOR_PROPERTIES_H_

#include <functional>
#include <string>

#include "node/node.h"
#include "node/property_store.h"

namespace btool::app::collector {

class Properties {
 public:
  static bool Local(const ::btool::node::PropertyStore *ps) {
    const bool *l = nullptr;
    ps->Read(kLocal, &l);
    return (l == nullptr ? false /* default */ : *l);
  }

  static std::string Root(const ::btool::node::PropertyStore *ps) {
    const std::string *r = nullptr;
    ps->Read(kRoot, &r);
    return (r == nullptr ? "" /* default */ : *r);
  }

  static void SetLocal(::btool::node::PropertyStore *ps, bool l) {
    ps->Write(kLocal, l);
  }

  static void SetRoot(::btool::node::PropertyStore *ps, std::string r) {
    ps->Write(kRoot, r);
  }

 private:
  static const char *kLocal;
  static const char *kRoot;
};

const std::vector<std::string> *ReadStringsProperty(
    const ::btool::node::PropertyStore *ps, const std::string &key);

void CollectStringsProperties(const ::btool::node::Node &n,
                              std::vector<std::string> *accumulator,
                              std::function<const std::vector<std::string> *(
                                  const ::btool::node::PropertyStore *)>
                                  f);

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_PROPERTIES_H_
