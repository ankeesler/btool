#include "app/collector/registry/registry_collectini.h"

#include "app/collector/base_collectini.h"
#include "app/collector/registry/registry.h"
#include "util/fs/fs.h"
#include "util/util.h"

namespace btool::app::collector::registry {

void RegistryCollectini::Collect(::btool::app::collector::Store *s) {
  Index i;
  std::string cache_dir = ::btool::util::Hex(r_->GetName());
  std::string key = ::btool::util::fs::Join(cache_dir, "index");
  if (!c_i_->Get(key, &i)) {
    r_->GetIndex(&i);
    c_i_->Set(key, i);
  }

  for (const auto &file : i.files) {
    Gaggle g;
    key = ::btool::util::fs::Join(cache_dir, file.sha256);
    if (!c_g_->Get(key, &g)) {
      r_->GetGaggle(file.path, &g);
      c_g_->Set(key, g);
    }

    std::string root = ::btool::util::fs::Join(cache_, file.sha256);
    gc_->Collect(s, &g, root);
  }
}

};  // namespace btool::app::collector::registry
