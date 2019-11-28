#include "app/cleaner/remove_aller_impl.h"

#include <string>

#include "util/fs/fs.h"

namespace btool::app::cleaner {

void RemoveAllerImpl::RemoveAll(const std::string &path) {
  ::btool::util::fs::RemoveAll(path);
}

};  // namespace btool::app::cleaner
