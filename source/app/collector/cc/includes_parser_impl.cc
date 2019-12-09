#include "app/collector/cc/includes_parser_impl.h"

#include <functional>
#include <string>

#include "util/fs/fs.h"

namespace btool::app::collector::cc {

void IncludesParserImpl::ParseIncludes(
    const std::string &path,
    std::function<void(const std::string &)> callback) {
  auto content = ::btool::util::fs::ReadFile(path);

  std::size_t index = 0;
  while (true) {
    index = content.find("#include", index);
    if (index == std::string::npos) {
      break;
    }

    std::size_t newline = content.find('\n', index);
    std::size_t start = content.find('"', index);
    if (newline < start) {
      index = newline + 1;
      continue;
    } else if (start == std::string::npos) {
      break;
    }

    std::size_t end = content.find('"', start + 1);
    if (end == std::string::npos) {
      break;
    }

    callback(content.substr(start + 1, end - start - 1));

    index = end + 1;
  }
}

};  // namespace btool::app::collector::cc
