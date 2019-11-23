#include "app/collector/registry/fs_registry.h"

#include <fstream>

#include "err.h"
#include "util/fs/fs.h"

namespace btool::app::collector::registry {

void FsRegistry::Initialize() {
  auto index_yml = ::btool::util::fs::Join(root_, "index.yml");
  std::ifstream index_ifs(index_yml);
  if (!index_ifs) {
    THROW_ERR("could not open " + index_yml);
  }
  s_->UnmarshalIndex(&index_ifs, &i_);

  for (const auto &f : i_.files) {
    auto gaggle_yml = ::btool::util::fs::Join(root_, f.path);
    std::ifstream gaggle_ifs(gaggle_yml);
    if (!gaggle_ifs) {
      THROW_ERR("could not open " + gaggle_yml);
    }

    Gaggle g;
    s_->UnmarshalGaggle(&gaggle_ifs, &g);
    gs_[f.path] = g;
  }

  initialized_ = true;
}

};  // namespace btool::app::collector::registry
