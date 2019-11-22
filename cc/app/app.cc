#include "app/app.h"

#include <string>

#include "err.h"
#include "log.h"
#include "node/node.h"

namespace btool::app {

bool App::Run(const std::string &target, bool clean, bool list, bool run,
              std::string *ret_err) {
  ::btool::node::Node *n = nullptr;
  std::string err;
  if (!collector_->Collect(target, &n, &err)) {
    *ret_err = ::btool::WrapErr(err, "collect");
    return true;
  }

  DEBUGS() << "collected graph from root " << n->name() << std::endl;

  if (clean) {
    if (!cleaner_->Clean(*n, &err)) {
      *ret_err = ::btool::WrapErr(err, "clean");
      return false;
    }
  } else if (list) {
    if (!lister_->List(*n, &err)) {
      *ret_err = ::btool::WrapErr(err, "list");
      return false;
    }
  } else {
    if (!builder_->Build(*n, &err)) {
      *ret_err = ::btool::WrapErr(err, "build");
      return false;
    }

    if (run) {
      if (!runner_->Run(*n, &err)) {
        *ret_err = ::btool::WrapErr(err, "run");
        return false;
      }
    }
  }

  return true;
}

};  // namespace btool::app
