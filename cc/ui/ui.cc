#include "ui/ui.h"

#include "core/log.h"

namespace btool::ui {

void UI::OnRun(const ::btool::node::Node &node) {
  INFO("running %s\n", node.Name().c_str());
}

};  // namespace btool::ui
