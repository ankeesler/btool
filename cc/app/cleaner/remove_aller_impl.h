#ifndef BTOOL_APP_CLEANER_REMOVEALLERIMPL_H_
#define BTOOL_APP_CLEANER_REMOVEALLERIMPL_H_

#include <string>

#include "cleaner.h"
#include "err.h"

namespace btool::app::cleaner {

class RemoveAllerImpl : public Cleaner::RemoveAller {
 public:
  void RemoveAll(const std::string &path) override;
};

};  // namespace btool::app::cleaner

#endif  // BTOOL_APP_CLEANER_REMOVEALLERIMPL_H_
