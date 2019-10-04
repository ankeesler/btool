#ifndef BTOOL_APP_RUNNER_RUNNER_H_
#define BTOOL_APP_RUNNER_RUNNER_H_

#include <string>

#include "node/node.h"

namespace btool::app::runner {

class Runner {
 public:
  bool Run(const ::btool::node::Node &node, std::string *err);
};

};  // namespace btool::app::runner

#endif  // BTOOL_APP_RUNNER_RUNNER_H_
