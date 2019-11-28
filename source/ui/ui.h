#ifndef BTOOL_UI_UI_H_
#define BTOOL_UI_UI_H_

#include "app/runner/runner.h"
#include "node/node.h"

namespace btool::ui {

class UI : public ::btool::app::runner::Runner::Callback {
 public:
  void OnRun(const ::btool::node::Node &node) override;
};

};  // namespace btool::ui

#endif  // BTOOL_UI_UI_H_
