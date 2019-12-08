#include "app/builder/currenter_impl.h"

#include <errno.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <sys/syscall.h>
#include <unistd.h>
#include <cstring>

#include <chrono>
#include <iostream>
#include <string>

#include "err.h"
#include "log.h"
#include "node/node.h"
#include "util/fs/fs.h"

namespace btool::app::builder {

bool CurrenterImpl::Current(const ::btool::node::Node &node) {
  if (!::btool::util::fs::Exists(node.name())) {
    DEBUGS() << node.name() << " does not exist, so it is not current"
             << std::endl;
    return false;
  }

  auto node_mod_time =
      ::btool::util::fs::ModTime<std::chrono::system_clock,
                                 std::chrono::nanoseconds>(node.name());

  std::chrono::time_point<std::chrono::system_clock, std::chrono::nanoseconds>
      latest_mod_time(std::chrono::nanoseconds(0));
  const ::btool::node::Node *latest_mod_time_node = nullptr;
  node.Visit([&](const ::btool::node::Node *dep) {
    auto dep_mod_time =
        ::btool::util::fs::ModTime<std::chrono::system_clock,
                                   std::chrono::nanoseconds>(dep->name());
    if (dep_mod_time > latest_mod_time) {
      latest_mod_time = dep_mod_time;
      latest_mod_time_node = dep;
    }
  });

  if (latest_mod_time_node != nullptr) {
    DEBUGS() << node.name() << " has latest modified ancestor "
             << latest_mod_time_node->name() << std::endl;
  }
  return latest_mod_time <= node_mod_time;
}

};  // namespace btool::app::builder
