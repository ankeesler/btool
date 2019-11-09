#ifndef BTOOL_APP_COLLECTOR_TRIVIALCOLLECTINI_H_
#define BTOOL_APP_COLLECTOR_TRIVIALCOLLECTINI_H_

#include <string>

#include "app/collector/base_collectini.h"
#include "app/collector/collector.h"
#include "app/collector/store.h"
#include "core/err.h"

namespace btool::app::collector {

class TrivialCollectini : public BaseCollectini {
 public:
  TrivialCollectini(const std::string &name) : name_(name) {}

  void Collect(Store *s) override {
    s->Put(name_);
    Notify(s, name_);
  }

 private:
  const std::string &name_;
};

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_TRIVIALCOLLECTINI_H_