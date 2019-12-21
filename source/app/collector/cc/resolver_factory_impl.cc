#include "app/collector/cc/resolver_factory_impl.h"

#include <memory>
#include <sstream>
#include <string>
#include <vector>

#include "app/collector/cc/properties.h"
#include "app/collector/properties.h"
#include "err.h"
#include "log.h"
#include "node/node.h"
#include "node/property_store.h"
#include "util/cmd.h"

namespace btool::app::collector::cc {

class CompileResolver : public ::btool::node::Node::Resolver {
 public:
  CompileResolver(std::string compiler, std::vector<std::string> flags)
      : compiler_(compiler), flags_(flags) {}

  void Resolve(const ::btool::node::Node &n) override {
    if (n.dependencies()->empty()) {
      THROW_ERR("node " + n.name() + "must have at least one dependency");
    }

    std::vector<std::string> include_dirs;
    ::btool::app::collector::CollectStringsProperties(
        n, &include_dirs,
        [](const ::btool::node::PropertyStore *ps)
            -> const std::vector<std::string> * {
          return Properties::IncludePaths(ps);
        });

    ::btool::util::Cmd cmd(compiler_.c_str());
    cmd.Arg("-o");
    cmd.Arg(n.name().c_str());
    cmd.Arg("-c");
    cmd.Arg(n.dependencies()->at(0)->name().c_str());
    for (const auto &id : include_dirs) {
      std::string flag("-I" + id);
      cmd.Arg(flag);
    }
    for (const auto &f : flags_) {
      cmd.Arg(f);
    }

    std::ostringstream out;
    std::ostringstream err;
    cmd.Stdout(&out);
    cmd.Stderr(&err);

    int ec = cmd.Run();

    DEBUGS() << "running linker invocation: " << cmd << std::endl;
    DEBUGS() << "out: " << out.str() << std::endl;
    DEBUGS() << "err: " << err.str() << std::endl;

    if (ec != 0) {
      THROW_ERR("compiler exit code = " + std::to_string(ec) +
                ", err: " + err.str());
    }
  }

 private:
  std::string compiler_;
  std::vector<std::string> flags_;
};

class ArchiveResolver : public ::btool::node::Node::Resolver {
 public:
  ArchiveResolver(std::string archiver) : archiver_(archiver) {}

  void Resolve(const ::btool::node::Node &n) override {
    ::btool::util::Cmd cmd(archiver_);
    cmd.Arg("rcs");
    cmd.Arg(n.name());
    for (const auto &d : *n.dependencies()) {
      cmd.Arg(d->name());
    }

    std::ostringstream out;
    std::ostringstream err;
    cmd.Stdout(&out);
    cmd.Stderr(&err);

    int ec = cmd.Run();

    DEBUGS() << "running archiver invocation: " << cmd << std::endl;
    DEBUGS() << "out: " << out.str() << std::endl;
    DEBUGS() << "err: " << err.str() << std::endl;

    if (ec != 0) {
      THROW_ERR("archiver exit code = " + std::to_string(ec) +
                ", err: " + err.str());
    }
  }

 private:
  std::string archiver_;
};

class LinkResolver : public ::btool::node::Node::Resolver {
 public:
  LinkResolver(std::string linker, std::vector<std::string> flags)
      : linker_(linker), flags_(flags) {}

  void Resolve(const ::btool::node::Node &n) override {
    std::vector<std::string> more_flags;
    ::btool::app::collector::CollectStringsProperties(
        n, &more_flags,
        [](const ::btool::node::PropertyStore *ps)
            -> const std::vector<std::string> * {
          return Properties::LinkFlags(ps);
        });

    ::btool::util::Cmd cmd(linker_.c_str());
    cmd.Arg("-o");
    cmd.Arg(n.name());
    for (const auto &d : *n.dependencies()) {
      cmd.Arg(d->name());
    }
    for (const auto &f : flags_) {
      cmd.Arg(f);
    }
    for (const auto &f : more_flags) {
      cmd.Arg(f);
    }

    std::ostringstream out;
    std::ostringstream err;
    cmd.Stdout(&out);
    cmd.Stderr(&err);

    int ec = cmd.Run();

    DEBUGS() << "running linker invocation: " << cmd << std::endl;
    DEBUG("out: %s\n", out.str().c_str());
    DEBUG("err: %s\n", err.str().c_str());

    if (ec != 0) {
      THROW_ERR("linker exit code = " + std::to_string(ec) +
                ", err: " + err.str());
    }
  }

 private:
  std::string linker_;
  std::vector<std::string> flags_;
};

::btool::node::Node::Resolver *ResolverFactoryImpl::NewCompileC() {
  auto cr = new CompileResolver(compiler_c_, compiler_c_flags_);
  resolvers_.emplace_back(cr);
  return cr;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewCompileCC() {
  auto cr = new CompileResolver(compiler_cc_, compiler_cc_flags_);
  resolvers_.emplace_back(cr);
  return cr;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewArchive() {
  auto ar = new ArchiveResolver(archiver_);
  resolvers_.emplace_back(ar);
  return ar;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewLinkC() {
  auto lr = new LinkResolver(linker_c_, linker_c_flags_);
  resolvers_.emplace_back(lr);
  return lr;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewLinkCC() {
  auto lr = new LinkResolver(linker_cc_, linker_cc_flags_);
  resolvers_.emplace_back(lr);
  return lr;
}

};  // namespace btool::app::collector::cc
