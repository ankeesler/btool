#include "app/builder/work_pool_impl.h"

#include <functional>
#include <iostream>
#include <thread>

#include "app/builder/channel.h"
#include "err.h"
#include "log.h"

namespace btool::app::builder {

WorkPoolImpl::WorkPoolImpl(std::size_t threads) {
  for (std::size_t i = 0; i < threads; ++i) {
    work_threads_.push_back(std::thread([this]() { WorkerLoop(); }));
  }
}

WorkPoolImpl::~WorkPoolImpl() {
  work_ch_.Close();
  done_ch_.Close();
  for (std::size_t i = 0; i < work_threads_.size(); ++i) {
    DEBUGS() << "joining thread " << i << std::endl;
    work_threads_[i].join();
  }
}

void WorkPoolImpl::Submit(std::function<const ::btool::node::Node *()> work) {
  if (!work_ch_.Tx(work)) {
    THROW_ERR("work channel has been closed!");
  }
}

const ::btool::node::Node *WorkPoolImpl::Receive(::btool::Err *err) {
  Done d;
  if (!done_ch_.Rx(&d)) {
    THROW_ERR("done channel has been closed!");
  }

  if (d.node == nullptr) {
    *err = d.err;
  }

  return d.node;
}

void WorkPoolImpl::WorkerLoop() {
  std::function<const ::btool::node::Node *()> work;
  while (work_ch_.Rx(&work)) {
    Done d = {
        .node = nullptr,
    };
    try {
      d.node = work();
    } catch (const ::btool::Err &err) {
      d.err = err;
    }

    if (!done_ch_.Tx(d)) {
      THROW_ERR("done channel has been closed!");
    }
  }
}

};  // namespace btool::app::builder
