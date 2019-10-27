#include "ui/ui.h"

#include "core/log.h"

namespace btool::ui {

void UI::OnRun(const ::btool::node::Node &node) {
  INFO("running %s\n", node.name().c_str());
}

};  // namespace btool::ui
