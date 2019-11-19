#ifndef BTOOL_APP_COLLECTOR_BASECOLLECTINI_H_
#define BTOOL_APP_COLLECTOR_BASECOLLECTINI_H_

#include <algorithm>
#include <string>
#include <vector>

#include "app/collector/collector.h"
#include "app/collector/store.h"

namespace btool::app::collector {

class BaseCollectini : public Collector::Collectini {
 public:
  BaseCollectini() { collectinis.push_back(this); }

  ~BaseCollectini() {
    auto it = std::find(collectinis.begin(), collectinis.end(), this);
    if (it != collectinis.end()) {
      collectinis.erase(it);
    }
  }

  virtual void Collect(Store *) override {}
  std::vector<std::string> Errors() override { return errors_; }

 protected:
  virtual void OnNotify(Store *s, const std::string &name) {}

  void AddError(std::string error) { errors_.push_back(error); }

  void Notify(Store *s, const std::string &name) {
    for (auto c : collectinis) {
      if (c != this) {
        c->OnNotify(s, name);
      }
    }
  }

 private:
  static std::vector<BaseCollectini *> collectinis;

  std::vector<std::string> errors_;
};

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_BASECOLLECTINI_H_
