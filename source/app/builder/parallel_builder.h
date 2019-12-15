#ifndef BTOOL_APP_BUILDER_PARALLEL_BUILDER_H_
#define BTOOL_APP_BUILDER_PARALLEL_BUILDER_H_

#include <functional>

#include "app/app.h"
#include "err.h"
#include "node/node.h"

namespace btool::app::builder {

class ParallelBuilder : public ::btool::app::App::Builder {
 public:
  class WorkPool {
   public:
    ~WorkPool() {}
    virtual void Submit(std::function<const ::btool::node::Node *()> work) = 0;
    virtual const ::btool::node::Node *Receive(::btool::Err *err) = 0;
  };

  class Currenter {
   public:
    virtual ~Currenter() {}
    virtual bool Current(const ::btool::node::Node &node) = 0;
  };

  // These callbacks will be called on the same thread.
  //
  // However, there many be overlapping calls to this Callback object from
  // different node resolutions.
  class Callback {
   public:
    virtual ~Callback() {}
    virtual void OnPreResolve(const ::btool::node::Node &node,
                              bool current) = 0;
    virtual void OnPostResolve(const ::btool::node::Node &node,
                               bool current) = 0;
  };

  ParallelBuilder(WorkPool *wp, Currenter *cu, Callback *ca)
      : wp_(wp), cu_(cu), ca_(ca) {}

  void Build(const ::btool::node::Node &node) override;

 private:
  const ::btool::node::Node *ReallyBuild(const ::btool::node::Node *node);

  WorkPool *wp_;
  Currenter *cu_;
  Callback *ca_;
};

};  // namespace btool::app::builder

#endif  // BTOOL_APP_BUILDER_PARALLEL_BUILDER_H_
