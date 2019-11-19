#include "collector.h"

#include <algorithm>
#include <iostream>
#include <sstream>
#include <string>
#include <vector>

#include "app/collector/store.h"
#include "err.h"
#include "log.h"
#include "node/node.h"

namespace btool::app::collector {

::btool::Err<::btool::node::Node *> Collector::Collect(
    const std::string &target) {
  for (auto c : cs_) {
    c->Collect(s_);

    auto errors = c->Errors();
    if (!errors.empty()) {
      std::stringstream ss{"collect errors:"};
      std::for_each(errors.begin(), errors.end(), [&ss](const std::string &s) {
        ss << std::endl << s;
      });
      return ::btool::Err<::btool::node::Node *>::Failure(ss.str().c_str());
    }
  }

  auto n = s_->Get(target);
  if (n == nullptr) {
    return ::btool::Err<::btool::node::Node *>::Failure("unknown target");
  }

  return ::btool::Err<::btool::node::Node *>::Success(n);
}

};  // namespace btool::app::collector
