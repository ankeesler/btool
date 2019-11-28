#include "app/runner/runner.h"

#include <iostream>
#include <string>

#include "err.h"
#include "util/cmd.h"

namespace btool::app::runner {

void Runner::Run(const ::btool::node::Node &node) {
  cb_->OnRun(node);

  ::btool::util::Cmd cmd(node.name().c_str());
  int ec = cmd.Run();
  if (ec != 0) {
    THROW_ERR("exit code = " + std::to_string(ec));
  }
}

};  // namespace btool::app::runner
