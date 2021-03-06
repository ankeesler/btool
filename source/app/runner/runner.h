#ifndef BTOOL_APP_RUNNER_RUNNER_H_
#define BTOOL_APP_RUNNER_RUNNER_H_

#include <string>

#include "app/app.h"
#include "err.h"
#include "node/node.h"

namespace btool::app::runner {

class Runner : public ::btool::app::App::Runner {
 public:
  class Callback {
   public:
    virtual ~Callback() {}
    virtual void OnRun(const ::btool::node::Node &node) = 0;
  };

  Runner(Callback *cb) : cb_(cb) {}

  void Run(const ::btool::node::Node &node) override;

 private:
  Callback *cb_;
};

};  // namespace btool::app::runner

#endif  // BTOOL_APP_RUNNER_RUNNER_H_
