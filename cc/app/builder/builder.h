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
    virtual ::btool::Err<bool> Current(
        const ::btool::node::Node &node) = 0;
  };

  Builder(Currenter *c) : c_(c) {}

  ::btool::VoidErr Build(const ::btool::node::Node &node) override;

 private:
  Currenter *c_;
};

};  // namespace btool::app::builder

#endif  // BTOOL_APP_BUILDER_BUILDER_H_
