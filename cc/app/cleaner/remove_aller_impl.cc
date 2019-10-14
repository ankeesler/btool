#include "app/cleaner/remove_aller_impl.h"

#include <string>

#include "util/fs/fs.h"

namespace btool::app::cleaner {

::btool::core::VoidErr RemoveAllerImpl::RemoveAll(const std::string &path) {
  return ::btool::util::fs::RemoveAll(path);
}

};  // namespace btool::app::cleaner
