#ifndef BTOOL_APP_COLLECTOR_BASECOLLECTINI_H_
#define BTOOL_APP_COLLECTOR_BASECOLLECTINI_H_

#include <vector>

#include "app/collector/collector.h"
#include "app/collector/store.h"

namespace btool::app::collector {

class BaseCollectini : public Collector::Collectini {
 public:
  BaseCollectini() { collectinis.push_back(this); }

  ~BaseCollectini() {}

  virtual void Collect(Store *s) = 0;

 protected:
  virtual void OnNotify(Store *s, const std::string &name) {}

  void Notify(Store *s, const std::string &name) {
    for (auto c : collectinis) {
      if (c != this) {
        c->OnNotify(s, name);
      }
    }
  }

 private:
  static std::vector<BaseCollectini *> collectinis;
};

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_BASECOLLECTINI_H_
