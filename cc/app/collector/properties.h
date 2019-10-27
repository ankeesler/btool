#ifndef BTOOL_APP_COLLECTOR_PROPERTIES_H_
#define BTOOL_APP_COLLECTOR_PROPERTIES_H_

#include <string>

#include "node/property_store.h"

namespace btool::app::collector {

class Properties {
 public:
  static bool Local(const ::btool::node::PropertyStore *ps) {
    const bool *l;
    ps->Read(kLocal, &l);
    return (l == nullptr ? false /* default */ : *l);
  }

  static void SetLocal(::btool::node::PropertyStore *ps, bool l) {
    ps->Write(kLocal, l);
  }

 private:
  static const std::string kLocal;
};

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_PROPERTIES_H_
