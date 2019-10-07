#ifndef BTOOL_APP_BUILDER_BUILDER_H_
#define BTOOL_APP_BUILDER_BUILDER_H_

#include <string>

#include "node/node.h"

namespace btool::app::builder {

class Builder {
 public:
  class Currenter {
   public:
    virtual ~Currenter() {}
    virtual bool Current(const ::btool::node::Node &node, bool *current,
                         std::string *err) = 0;
  };

  Builder(Currenter *c) : c_(c) {}

  bool Build(const ::btool::node::Node &node, std::string *err);

 private:
  Currenter *c_;
};

};  // namespace btool::app::builder

#endif  // BTOOL_APP_BUILDER_BUILDER_H_
