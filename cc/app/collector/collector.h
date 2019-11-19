#ifndef BTOOL_APP_COLLECTOR_COLLECTOR_H_
#define BTOOL_APP_COLLECTOR_COLLECTOR_H_

#include <string>
#include <vector>

#include "app/app.h"
#include "app/collector/store.h"
#include "err.h"
#include "node/node.h"

namespace btool::app::collector {

class Collector : public ::btool::app::App::Collector {
 public:
  class Collectini {
   public:
    virtual ~Collectini() {}
    virtual void Collect(Store *) = 0;
    virtual std::vector<std::string> Errors() = 0;
  };

  Collector(Store *s) : s_(s) {}

  void AddCollectini(Collectini *c) { cs_.push_back(c); }
  ::btool::Err<::btool::node::Node *> Collect(
      const std::string &target) override;

 private:
  Store *s_;
  std::vector<Collectini *> cs_;
};

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_COLLECTOR_H_
