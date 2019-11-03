#include "app/collector/cc/includes_parser_impl.h"

#include <functional>
#include <string>

#include "core/err.h"
#include "util/fs/fs.h"

namespace btool::app::collector::cc {

::btool::core::VoidErr IncludesParserImpl::ParseIncludes(
    const std::string &path,
    std::function<void(const std::string &)> callback) {
  auto err = ::btool::util::fs::ReadFile(path);
  if (err) {
    return ::btool::core::VoidErr::Failure(err.Msg());
  }
  auto content = err.Ret();

  std::size_t index = 0;
  while (true) {
    index = content.find("#include", index);
    if (index == std::string::npos) {
      break;
    }

    size_t start = content.find('"', index);
    if (start == std::string::npos) {
      break;
    }

    size_t end = content.find('"', start + 1);
    if (end == std::string::npos) {
      break;
    }

    callback(content.substr(start + 1, end - start - 1));

    index = end + 1;
  }

  return ::btool::core::VoidErr::Success();
}

};  // namespace btool::app::collector::cc
