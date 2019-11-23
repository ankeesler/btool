#include "app/collector/registry/http_registry.h"

#include <fstream>
#include <functional>
#include <iostream>
#include <string>

#include "err.h"
#include "log.h"
#include "util/download.h"
#include "util/fs/fs.h"

namespace btool::app::collector::registry {

static void Get(const std::string &url,
                std::function<void(std::istream *)> unmarshal);

void HttpRegistry::GetIndex(Index *i) {
  Get(url_, [this, i](std::istream *is) { s_->UnmarshalIndex(is, i); });
}

void HttpRegistry::GetGaggle(std::string name, Gaggle *g) {
  Get(url_ + "/" + name,
      [this, g](std::istream *is) { s_->UnmarshalGaggle(is, g); });
}

static void Get(const std::string &url,
                std::function<void(std::istream *)> unmarshal) {
  auto dir = ::btool::util::fs::TempDir();
  auto file = ::btool::util::fs::Join(dir, "file");

  ::btool::util::Download(url, file);
  DEBUGS() << "dowloaded file " << ::btool::util::fs::ReadFile(file)
           << std::endl;

  std::ifstream ifs(file);
  if (!ifs) {
    THROW_ERR("cannot open downloaded file " + file);
  }

  unmarshal(&ifs);

  ::btool::util::fs::RemoveAll(dir);
}

};  // namespace btool::app::collector::registry
