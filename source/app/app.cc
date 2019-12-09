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

static std::chrono::system_clock::duration Time(std::function<void(void)> f);

static std::string DurationString(std::chrono::system_clock::duration d);

void App::Run(const std::string &target, bool clean, bool list, bool run) {
  std::chrono::system_clock::duration collect_time;
  std::chrono::system_clock::duration clean_time;
  std::chrono::system_clock::duration list_time;
  std::chrono::system_clock::duration build_time;
  std::chrono::system_clock::duration run_time;

  auto collect_start = std::chrono::system_clock::now();
  auto n = collector_->Collect(target);
  auto collect_end = std::chrono::system_clock::now();
  collect_time = collect_end - collect_start;

  if (clean) {
    clean_time = Time([this, n] { cleaner_->Clean(*n); });
  } else if (list) {
    list_time = Time([this, n] { lister_->List(*n); });
  } else {
    build_time = Time([this, n] { builder_->Build(*n); });

    if (run) {
      run_time = Time([this, n] { runner_->Run(*n); });
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

static std::chrono::system_clock::duration Time(std::function<void(void)> f) {
  auto start = std::chrono::system_clock::now();
  f();
  auto end = std::chrono::system_clock::now();
  return end - start;
}

static std::string DurationString(std::chrono::system_clock::duration d) {
  std::stringstream ss;
  ss << ::btool::util::CommaSeparatedNumber(
            std::chrono::duration_cast<std::chrono::milliseconds>(d).count())
     << " ms";
  return ss.str();
}

};  // namespace btool::app
