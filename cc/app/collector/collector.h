#ifndef BTOOL_APP_COLLECTOR_COLLECTOR_H_
#define BTOOL_APP_COLLECTOR_COLLECTOR_H_

#include <vector>

#include "app/app.h"
#include "node/store.h"

namespace btool::app::collector {

class Collector : public ::btool::app::App::Collector {
 public:
  class Collectini {
   public:
    virtual ~Collectini() {}
    virtual ::btool::core::VoidErr Collect(::btool::node::Store *) = 0;
  };

  Collector(std::vector<Collectini *> *cs) : cs_(cs) {}

  ::btool::core::VoidErr Collect() override;

 private:
  std::vector<Collectini *> *cs_;
};

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_COLLECTOR_H_
