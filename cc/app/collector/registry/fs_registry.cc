#include "app/collector/registry/fs_registry.h"

#include <fstream>

#include "err.h"
#include "util/fs/fs.h"

namespace btool::app::collector::registry {

void FsRegistry::Initialize() {
  auto index_yml = ::btool::util::fs::Join(root_, "index.yml");
  std::ifstream ifs(index_yml);
  if (!ifs) {
    THROW_ERR("could not open " + index_yml);
  }
  s_->UnmarshalIndex(&ifs, &i_);

  initialized_ = true;
}

};  // namespace btool::app::collector::registry
