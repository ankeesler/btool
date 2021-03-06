#ifndef BTOOL_APP_BUILDER_CURRENTERIMPL_H_
#define BTOOL_APP_BUILDER_CURRENTERIMPL_H_

#include "app/builder/parallel_builder.h"
#include "err.h"
#include "node/node.h"

namespace btool::app::builder {

class CurrenterImpl : public ParallelBuilder::Currenter {
 public:
  bool Current(const ::btool::node::Node &node) override;
};

};  // namespace btool::app::builder

#endif  // BTOOL_APP_BUILDER_CURRENTERIMPL_H_
