#ifndef BTOOL_APP_COLLECTOR_PROPERTIES_H_
#define BTOOL_APP_COLLECTOR_PROPERTIES_H_

#include <string>

#include "node/property_store.h"

namespace btool::app::collector {

class Properties {
 public:
  Properties(::btool::node::PropertyStore *ps) : ps_(ps) {}

  bool local() const {
    const bool *l;
    ps_->Read(kLocal, &l);
    return (l == nullptr ? false /* default */ : *l);
  }

  void set_local(bool l) { ps_->Write(kLocal, l); }

 private:
  const std::string kLocal = "io.btool.app.collector.local";

  ::btool::node::PropertyStore *ps_;
};

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_PROPERTIES_H_
