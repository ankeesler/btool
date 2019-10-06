#include "app/runner/runner.h"

#include <iostream>

#include "core/cmd.h"

namespace btool::app::runner {

bool Runner::Run(const ::btool::node::Node &node) {
  cb_->OnRun(node);

  ::btool::core::Cmd cmd(node.Name().c_str());
  return (cmd.Run() == 0);
}

};  // namespace btool::app::runner
