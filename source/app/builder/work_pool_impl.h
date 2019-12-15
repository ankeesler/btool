#ifndef BTOOL_APP_BUILDER_WORK_POOL_IMPL_H_
#define BTOOL_APP_BUILDER_WORK_POOL_IMPL_H_

#include <functional>
#include <thread>
#include <vector>

#include "app/builder/channel.h"
#include "app/builder/parallel_builder.h"
#include "err.h"
#include "node/node.h"

namespace btool::app::builder {

class WorkPoolImpl : public ParallelBuilder::WorkPool {
 public:
  WorkPoolImpl(std::size_t threads);
  ~WorkPoolImpl();

  void Submit(std::function<const ::btool::node::Node *()> work) override;
  const ::btool::node::Node *Receive(::btool::Err *err) override;

 private:
  struct Done {
    const ::btool::node::Node *node;
    ::btool::Err err;
  };

  std::vector<std::thread> work_threads_;

  Channel<std::function<const ::btool::node::Node *()>> work_ch_;
  Channel<Done> done_ch_;

  void WorkerLoop();
};

};  // namespace btool::app::builder

#endif  // BTOOL_APP_BUILDER_WORK_POOL_IMPL_H_
