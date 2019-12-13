#include "app/app.h"

#include <chrono>
#include <functional>
#include <iostream>
#include <ostream>
#include <sstream>
#include <string>

#include "err.h"
#include "log.h"
#include "node/node.h"
#include "util/util.h"

namespace btool::app {

static std::string DurationString(std::chrono::steady_clock::duration d);

void App::Run(const std::string &target, bool clean, bool list, bool run) {
  std::chrono::steady_clock::duration collect_time;
  std::chrono::steady_clock::duration clean_time;
  std::chrono::steady_clock::duration list_time;
  std::chrono::steady_clock::duration build_time;
  std::chrono::steady_clock::duration run_time;

  ::btool::node::Node const *n = nullptr;
  collect_time = ::btool::util::Time(
      [this, target, &n]() { n = collector_->Collect(target); });

  if (clean) {
    clean_time = ::btool::util::Time([this, n] { cleaner_->Clean(*n); });
  } else if (list) {
    list_time = ::btool::util::Time([this, n] { lister_->List(*n); });
  } else {
    build_time = ::btool::util::Time([this, n] { builder_->Build(*n); });

    if (run) {
      run_time = ::btool::util::Time([this, n] { runner_->Run(*n); });
    }
  }

  std::ostream &os = INFOS() << "summary:" << std::endl;
  os << "  collect: " << DurationString(collect_time) << std::endl;
  if (clean) {
    os << "  clean: " << DurationString(clean_time) << std::endl;
  }
  if (list) {
    os << "  list: " << DurationString(list_time) << std::endl;
  }
  if (!clean && !list) {
    os << "  build: " << DurationString(build_time) << std::endl;
  }
  if (run) {
    os << "  run: " << DurationString(run_time) << std::endl;
  }
}

static std::string DurationString(std::chrono::steady_clock::duration d) {
  std::stringstream ss;
  ss << ::btool::util::CommaSeparatedNumber(
            std::chrono::duration_cast<std::chrono::milliseconds>(d).count())
     << " ms";
  return ss.str();
}

};  // namespace btool::app
