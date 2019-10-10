#include "app/app.h"

#include "core/err.h"
#include "node/node.h"

namespace btool::app {

::btool::core::VoidErr App::Run(bool clean, bool list, bool run) {
  auto err = collector_->Collect();
  if (err) {
    return err;
  }

  ::btool::node::Node n("a");

  if (clean) {
    err = cleaner_->Clean(n);
    if (err) {
      return err;
    }
  } else if (list) {
    err = lister_->List(n);
    if (err) {
      return err;
    }
  } else {
    err = builder_->Build(n);
    if (err) {
      return err;
    }

    if (run) {
      err = runner_->Run(n);
      if (err) {
        return err;
      }
    }
  }

  return ::btool::core::VoidErr::Success();
}

};  // namespace btool::app
