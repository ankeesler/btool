#include "app/runner/runner.h"

#include <iostream>

#include "core/err.h"
#include "util/cmd.h"

namespace btool::app::runner {

::btool::core::VoidErr Runner::Run(const ::btool::node::Node &node) {
  cb_->OnRun(node);

  ::btool::util::Cmd cmd(node.Name().c_str());
  return (cmd.Run() == 0 ? ::btool::core::VoidErr::Success()
                         : ::btool::core::VoidErr::Failure("exit code != 0"));
}

};  // namespace btool::app::runner
