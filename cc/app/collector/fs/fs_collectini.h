#ifndef BTOOL_APP_COLLECTOR_FS_FS_H_
#define BTOOL_APP_COLLECTOR_FS_FS_H_

#include "app/collector/collector.h"
#include "node/node.h"

namespace btool::app::collector::fs {

class FSCollectini : public ::btool::app::collector::Collector::Collectini {
 public:
  FSCollectini(const char *root) : root_(root) {}

  ::btool::core::VoidErr Collect(::btool::node::Store *) override;

 private:
  const char *root_;
};

};  // namespace btool::app::collector::fs

#endif  // BTOOL_APP_COLLECTOR_FS_FS_H_
