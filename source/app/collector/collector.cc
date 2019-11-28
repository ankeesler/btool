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

::btool::node::Node *Collector::Collect(const std::string &target) {
  for (auto c : cs_) {
    c->Collect(s_);

    auto errors = c->Errors();
    if (!errors.empty()) {
      std::stringstream ss{"collect errors:"};
      std::for_each(errors.begin(), errors.end(), [&ss](const std::string &s) {
        ss << std::endl << s;
      });
      THROW_ERR(ss.str());
    }
  }

  auto n = s_->Get(target);
  if (n == nullptr) {
    THROW_ERR("unknown target");
  }

  return n;
}

};  // namespace btool::app::collector
