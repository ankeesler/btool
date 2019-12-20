#include "app/collector/cc/includes_parser_impl.h"

#include <functional>
#include <string>

#include "util/fs/fs.h"

namespace btool::app::collector::cc {

void IncludesParserImpl::ParseIncludes(
    const std::string &path,
    std::function<void(const std::string &, bool)> callback) {
  auto content = ::btool::util::fs::ReadFile(path);

  std::size_t index = 0;
  while (true) {
    index = content.find("#include", index);
    if (index == std::string::npos) {
      break;
    }

    std::size_t newline = content.find('\n', index);
    std::size_t system_start = content.find('<', index);
    std::size_t local_start = content.find('"', index);

    std::size_t start;
    char char_end;
    if (newline < system_start && newline < local_start) {
      index = newline + 1;
      continue;
    } else if (system_start == std::string::npos &&
               local_start == std::string::npos) {
      break;
    } else if (system_start == std::string::npos) {
      start = local_start;
      char_end = '"';
    } else if (local_start == std::string::npos) {
      start = system_start;
      char_end = '>';
    } else if (system_start < local_start) {
      start = system_start;
      char_end = '>';
    } else {  // system_start >= local_start
      start = local_start;
      char_end = '"';
    }

    std::size_t end = content.find(char_end, start + 1);
    if (end == std::string::npos) {
      break;
    }

    callback(content.substr(start + 1, end - start - 1), start == system_start);

    index = end + 1;
  }
}

};  // namespace btool::app::collector::cc
