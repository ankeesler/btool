#include "ui/ui.h"

#include <chrono>
#include <iostream>
#include <mutex>
#include <ostream>
#include <string>

#include "log.h"
#include "util/fs/fs.h"
#include "util/string/string.h"
#include "util/util.h"

namespace btool::ui {

void UI::OnRun(const ::btool::node::Node &node) {
  INFO("running %s\n", node.name().c_str());
}

void UI::OnPreResolve(const ::btool::node::Node &node, bool current) {
  resolve_start_[node.name()] = std::chrono::system_clock::now();

  std::unique_lock<std::mutex> lock(mtx_);
  INFOS() << "resolving " << MakeNamePretty(node.name(), cache_) << std::endl;
}

void UI::OnPostResolve(const ::btool::node::Node &node, bool current) {
  auto resolve_end = std::chrono::system_clock::now();

  std::unique_lock<std::mutex> lock(mtx_);
  std::ostream &os = INFOS()
                     << "resolving " << MakeNamePretty(node.name(), cache_);
  if (!current) {
    auto resolve_duration =
        std::chrono::duration_cast<std::chrono::milliseconds>(
            resolve_end - resolve_start_[node.name()]);
    os << " (" << ::btool::util::CommaSeparatedNumber(resolve_duration.count())
       << " ms)";
  } else {
    os << " (up to date)";
  }
  os << std::endl;
}

std::string MakeNamePretty(const std::string &name, const std::string &cache) {
  if (::btool::util::string::HasPrefix(name, cache)) {
    // This should come out looking like:
    //   $CACHE/abc.../tuna.o
    //          ^      ^
    //   start of sha  ^
    //                 ^
    //         basename of the file
    return "$CACHE" + name.substr(cache.size(), 4) + ".../" +
           ::btool::util::fs::Base(name);
  }

  const std::string line_prefix = "btool | info | resolving ";
  const std::string line_suffix = " (up to date)";
  const std::size_t max_line_size = 80;
  const std::size_t max_file_size =
      max_line_size - line_prefix.size() - line_suffix.size();
  if (name.size() > max_file_size) {
    // This should come out looking like:
    //   source/app.../app.cc
    //          ^      ^
    //          ^    basename of the file
    //   start of sha (for cache file case below), 10 characters over
    return name.substr(0, 10) + ".../" + ::btool::util::fs::Base(name);
  }

  return name;
}
};  // namespace btool::ui
