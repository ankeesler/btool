#include "app/app.h"

#include <string>

#include "err.h"
#include "log.h"
#include "node/node.h"

namespace btool::app {

void App::Run(const std::string &target, bool clean, bool list, bool run) {
  auto n = collector_->Collect(target);
  DEBUGS() << "collected graph from root " << n->name() << std::endl;

  if (clean) {
    cleaner_->Clean(*n);
  } else if (list) {
    lister_->List(*n);
  } else {
    builder_->Build(*n);

    if (run) {
      runner_->Run(*n);
    }
  }
}

};  // namespace btool::app
