#ifndef BTOOL_UI_UI_H_
#define BTOOL_UI_UI_H_

#include <chrono>
#include <mutex>
#include <string>

#include "app/builder/parallel_builder.h"
#include "app/runner/runner.h"
#include "node/node.h"

namespace btool::ui {

class UI : public ::btool::app::runner::Runner::Callback,
           public ::btool::app::builder::ParallelBuilder::Callback {
 public:
  UI(std::string cache) : cache_(cache) {}

  void OnRun(const ::btool::node::Node &node) override;

  void OnPreResolve(const ::btool::node::Node &node, bool current) override;
  void OnPostResolve(const ::btool::node::Node &node, bool current) override;

 private:
  std::string cache_;

  std::map<std::string, std::chrono::time_point<std::chrono::system_clock>>
      resolve_start_;
  std::mutex mtx_;
};

std::string MakeNamePretty(const std::string &name, const std::string &cache);

};  // namespace btool::ui

#endif  // BTOOL_UI_UI_H_
