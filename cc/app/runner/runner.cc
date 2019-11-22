#include "app/runner/runner.h"

#include <iostream>
#include <string>

#include "err.h"
#include "util/cmd.h"

namespace btool::app::runner {

bool Runner::Run(const ::btool::node::Node &node, std::string *ret_err) {
  cb_->OnRun(node);

  ::btool::util::Cmd cmd(node.name().c_str());
  int ec = cmd.Run();
  if (ec != 0) {
    *ret_err = ::btool::WrapErr("exit code " + ec, "run");
    return false;
  } else {
    return true;
  }
}

};  // namespace btool::app::runner
