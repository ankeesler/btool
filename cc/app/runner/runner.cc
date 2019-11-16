#include "app/runner/runner.h"

#include <iostream>

#include "err.h"
#include "util/cmd.h"

namespace btool::app::runner {

::btool::VoidErr Runner::Run(const ::btool::node::Node &node) {
  cb_->OnRun(node);

  ::btool::util::Cmd cmd(node.name().c_str());
  return (cmd.Run() == 0 ? ::btool::VoidErr::Success()
                         : ::btool::VoidErr::Failure("exit code != 0"));
}

};  // namespace btool::app::runner
