#include "ui/ui.h"

#include <iostream>
#include <ostream>
#include <string>

#include "log.h"
#include "util/fs/fs.h"
#include "util/string/string.h"

namespace btool::ui {

void UI::OnRun(const ::btool::node::Node &node) {
  INFO("running %s\n", node.name().c_str());
}

void UI::OnResolve(const ::btool::node::Node &node, bool current) {
  std::ostream &os = INFOS()
                     << "resolving " << MakeNamePretty(node.name(), cache_);
  if (current) {
    os << " (up to date)";
  }
  os << std::endl;
}

std::string MakeNamePretty(const std::string &name, const std::string &cache) {
  if (!::btool::util::string::HasPrefix(name, cache)) {
    return name;
  }

  // This should come out looking like:
  //   $CACHE/abc.../tuna.o
  //          ^      ^
  //   start of sha  ^
  //                 ^
  //         basename of the file
  return "$CACHE" + name.substr(cache.size(), 4) + ".../" +
         ::btool::util::fs::Base(name);
}

};  // namespace btool::ui
