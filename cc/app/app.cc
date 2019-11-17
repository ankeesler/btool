#include "app/app.h"

#include <string>

#include "err.h"
#include "log.h"
#include "node/node.h"

namespace btool::app {

::btool::VoidErr App::Run(const std::string &target, bool clean, bool list,
                          bool run) {
  auto collect_err = collector_->Collect(target);
  if (collect_err) {
    return ::btool::VoidErr::Failure(collect_err.Msg());
  }

  ::btool::node::Node *n = collect_err.Ret();
  DEBUG("collected graph from root %s\n", n->name().c_str());

  ::btool::VoidErr err;
  if (clean) {
    err = cleaner_->Clean(*n);
    if (err) {
      return err;
    }
  } else if (list) {
    err = lister_->List(*n);
    if (err) {
      return err;
    }
  } else {
    err = builder_->Build(*n);
    if (err) {
      return err;
    }

    if (run) {
      err = runner_->Run(*n);
      if (err) {
        return err;
      }
    }
  }

  return err;
}

};  // namespace btool::app
