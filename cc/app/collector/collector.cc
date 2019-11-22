#include "collector.h"

#include <algorithm>
#include <iostream>
#include <sstream>
#include <string>
#include <vector>

#include "app/collector/store.h"
#include "log.h"
#include "node/node.h"

namespace btool::app::collector {

bool Collector::Collect(const std::string &target, ::btool::node::Node **ret_n,
                        std::string *ret_err) {
  for (auto c : cs_) {
    c->Collect(s_);

    auto errors = c->Errors();
    if (!errors.empty()) {
      *ret_err = "collect errors:";
      std::for_each(errors.begin(), errors.end(),
                    [ret_err](const std::string &s) { *ret_err += "\n" + s; });
      return false;
    }
  }

  *ret_n = s_->Get(target);
  if (*ret_n == nullptr) {
    *ret_err = "unknown target";
    return false;
  }

  return true;
}

};  // namespace btool::app::collector
