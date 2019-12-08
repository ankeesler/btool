#include "app/collector/registry/resolver_factory_impl.h"

#include <fstream>
#include <string>
#include <vector>

#include "app/collector/registry/registry.h"
#include "err.h"
#include "log.h"
#include "node/node.h"
#include "util/cmd.h"
#include "util/download.h"
#include "util/fs/fs.h"
#include "util/sha256.h"

namespace btool::app::collector::registry {

class DownloadResolver : public ::btool::node::Node::Resolver {
 public:
  DownloadResolver(std::string url, std::string sha256)
      : url_(url), sha256_(sha256) {}

  void Resolve(const ::btool::node::Node &n) override {
    ::btool::util::fs::MkdirAll(::btool::util::fs::Dir(n.name()));
    ::btool::util::Download(url_, n.name());
    std::ifstream ifs(n.name());
    std::string actual_sha256 = ::btool::util::SHA256(&ifs);
    if (sha256_ != actual_sha256) {
      THROW_ERR("sha256 mismatch for node " + n.name() + " ex=" + sha256_ +
                " != ac=" + actual_sha256);
    }
  }

 private:
  std::string url_;
  std::string sha256_;
};

class UnzipResolver : public ::btool::node::Node::Resolver {
 public:
  void Resolve(const ::btool::node::Node &n) override {
    // if (n.dependencies()->empty()) {
    //  THROW_ERR("unzip resolve target " + n.name() +
    //            " must have at least one dependency");
    //}
    // auto zipfile = n.dependencies()->at(0)->name();
    // auto dir = ::btool::util::fs::Dir(zipfile);
    //::btool::util::Unzip(zipfile, dir);
  }
};

class UntarResolver : public ::btool::node::Node::Resolver {
 public:
  void Resolve(const ::btool::node::Node &n) override {
    if (n.dependencies()->empty()) {
      THROW_ERR("unzip resolve target " + n.name() +
                " must have at least one dependency");
    }
    auto tarfile = n.dependencies()->at(0)->name();
    auto dir = ::btool::util::fs::Dir(tarfile);
    ::btool::util::Cmd cmd("tar");
    cmd.Arg("mxzf");
    cmd.Arg(tarfile);
    cmd.Arg("-C");
    cmd.Arg(dir);

    std::ostringstream out;
    std::ostringstream err;
    cmd.Stdout(&out);
    cmd.Stderr(&err);

    int ec = cmd.Run();

    DEBUGS() << "tar out: " << out.str() << std::endl;
    DEBUGS() << "tar err: " << err.str() << std::endl;

    if (ec != 0) {
      THROW_ERR("tar exit code = " + std::to_string(ec) +
                ", err: " + err.str());
    }
  }
};

class CmdResolver : public ::btool::node::Node::Resolver {
 public:
  CmdResolver(std::string path, std::vector<std::string> args, std::string dir)
      : path_(path), args_(args), dir_(dir) {}

  void Resolve(const ::btool::node::Node &n) override {
    ::btool::util::Cmd cmd(path_);
    for (const auto &arg : args_) {
      cmd.Arg(arg);
    }
    cmd.Dir(dir_);

    std::ostringstream out;
    std::ostringstream err;
    cmd.Stdout(&out);
    cmd.Stderr(&err);

    int ec = cmd.Run();

    DEBUGS() << path_ << " dir: " << cmd.Dir() << std::endl;
    DEBUGS() << path_ << " out: " << out.str() << std::endl;
    DEBUGS() << path_ << " err: " << err.str() << std::endl;

    if (ec != 0) {
      THROW_ERR(path_ + " exit code = " + std::to_string(ec) +
                ", err: " + err.str());
    }
  }

 private:
  std::string path_;
  std::vector<std::string> args_;
  std::string dir_;
};

::btool::node::Node::Resolver *ResolverFactoryImpl::NewDownload(
    const std::string &url, const std::string &sha256) {
  auto r = new DownloadResolver(url, sha256);
  allocations_.push_back(r);
  return r;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewUnzip() {
  auto r = new UnzipResolver();
  allocations_.push_back(r);
  return r;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewUntar() {
  auto r = new UntarResolver();
  allocations_.push_back(r);
  return r;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewCmd(
    const std::string &path, const std::vector<std::string> &args,
    const std::string &dir) {
  auto r = new CmdResolver(path, args, dir);
  allocations_.push_back(r);
  return r;
}

};  // namespace btool::app::collector::registry
