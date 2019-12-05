#ifndef BTOOL_APP_BUILDER_BUILDER_H_
#define BTOOL_APP_BUILDER_BUILDER_H_

#include <string>

#include "app/app.h"
#include "err.h"
#include "node/node.h"

namespace btool::app::builder {

class Builder : public ::btool::app::App::Builder {
 public:
  class Currenter {
   public:
    virtual ~Currenter() {}
    virtual bool Current(const ::btool::node::Node &node) = 0;
  };

  class Callback {
   public:
    virtual ~Callback() {}
    virtual void OnResolve(const ::btool::node::Node &node, bool current) = 0;
  };

  Builder(Currenter *cu, Callback *ca) : cu_(cu), ca_(ca) {}

  void Build(const ::btool::node::Node &node) override;

 private:
  Currenter *cu_;
  Callback *ca_;
};

};  // namespace btool::app::builder

#endif  // BTOOL_APP_BUILDER_BUILDER_H_
