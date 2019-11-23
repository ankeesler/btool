#include "app/collector/registry/registry_collectini.h"

#include "app/collector/base_collectini.h"
#include "app/collector/registry/registry.h"
#include "util/fs/fs.h"

namespace btool::app::collector::registry {

void RegistryCollectini::Collect(::btool::app::collector::Store *s) {
  Index i;
  r_->GetIndex(&i);

  for (const auto &file : i.files) {
    Gaggle g;
    r_->GetGaggle(file.path, &g);

    std::string root = ::btool::util::fs::Join(cache_, file.sha256);
    gc_->Collect(s, &g, root);
  }
}

};  // namespace btool::app::collector::registry
