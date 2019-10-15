#ifndef BTOOL_APP_COLLECTOR_FS_FS_H_
#define BTOOL_APP_COLLECTOR_FS_FS_H_

#include <string>

#include "app/collector/collector.h"
#include "app/collector/store.h"
#include "node/node.h"

namespace btool::app::collector::fs {

class FSCollectini : public ::btool::app::collector::Collector::Collectini {
 public:
  FSCollectini(const std::string &root) : root_(root) {}

  ::btool::core::VoidErr Collect(::btool::app::collector::Store *) override;

 private:
  std::string root_;
};

};  // namespace btool::app::collector::fs

#endif  // BTOOL_APP_COLLECTOR_FS_FS_H_
