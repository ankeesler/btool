#include "app/collector/registry/resolver_factory_impl.h"

#include <string>

#include "app/collector/registry/registry.h"
#include "err.h"
#include "node/node.h"
#include "util/download.h"
#include "util/fs/fs.h"
#include "util/unzip_btool.h"

namespace btool::app::collector::registry {

class DownloadResolver : public ::btool::node::Node::Resolver {
 public:
  DownloadResolver(std::string url, std::string sha256)
      : url_(url), sha256_(sha256) {}

  void Resolve(const ::btool::node::Node &n) override {
    ::btool::util::Download(url_, n.name());
    (void)sha256_;  // TODO: check sha256!
  }

 private:
  std::string url_;
  std::string sha256_;
};

class UnzipResolver : public ::btool::node::Node::Resolver {
 public:
  void Resolve(const ::btool::node::Node &n) override {
    if (n.dependencies()->empty()) {
      THROW_ERR("unzip resolve target " + n.name() +
                " must have at least one dependency");
    }
    auto zipfile = n->dependencies()->at(0)->name();
    auto dir = ::btool::util::fs::Dir(zipfile);
    ::btool::util::Unzip(zipfile, dir);
  }
};

::btool::node::Node::Resolver *ResolverFactoryImpl::NewDownload(
    const std::string &url, const std::string &sha256) {
  auto r = new DownloadResolver(url, sha256);
  allocations_.push_back(r);
  return r;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewUnzip(
    const std::string &output_dir) {
  return nullptr;
}

};  // namespace btool::app::collector::registry
