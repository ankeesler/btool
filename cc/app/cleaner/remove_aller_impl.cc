#include "app/cleaner/remove_aller_impl.h"

#include <string>

namespace btool::app::cleaner {

::btool::core::VoidErr RemoveAllerImpl::RemoveAll(const std::string &path) {
  return ::btool::core::VoidErr::Success();
}

};  // namespace btool::app::cleaner
