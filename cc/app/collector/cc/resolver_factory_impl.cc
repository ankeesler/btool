#include "app/collector/cc/resolver_factory_impl.h"

#include <sstream>
#include <string>
#include <vector>

#include "core/err.h"
#include "core/log.h"
#include "node/node.h"
#include "util/cmd.h"

namespace btool::app::collector::cc {

class CompileResolver : public ::btool::node::Node::Resolver {
 public:
  CompileResolver(std::string compiler, std::vector<std::string> flags,
                  std::vector<std::string> include_dirs,
                  std::vector<std::string> more_flags)
      : compiler_(compiler),
        flags_(flags),
        include_dirs_(include_dirs),
        more_flags_(more_flags) {}

  ::btool::core::VoidErr Resolve(const ::btool::node::Node &n) override {
    if (n.dependencies()->empty()) {
      return ::btool::core::VoidErr::Failure(
          "must have at least one dependency");
    }

    ::btool::util::Cmd cmd(compiler_.c_str());
    cmd.Arg("-o");
    cmd.Arg(n.name().c_str());
    cmd.Arg("-c");
    cmd.Arg(n.dependencies()->at(0)->name().c_str());
    for (const auto &id : include_dirs_) {
      std::string flag("-I" + id);
      cmd.Arg(flag);
    }
    for (const auto &f : flags_) {
      cmd.Arg(f);
    }
    for (const auto &f : more_flags_) {
      cmd.Arg(f);
    }

    std::ostringstream out;
    std::ostringstream err;
    cmd.Stdout(&out);
    cmd.Stderr(&err);

    int ec = cmd.Run();

    DEBUG("out: %s\n", out.str().c_str());
    DEBUG("err: %s\n", err.str().c_str());

    return (ec == 0
                ? ::btool::core::VoidErr::Success()
                : ::btool::core::VoidErr::Failure("failed to run compiler"));
  }

 private:
  std::string compiler_;
  std::vector<std::string> flags_;
  std::vector<std::string> include_dirs_;
  std::vector<std::string> more_flags_;
};

class LinkResolver : public ::btool::node::Node::Resolver {
 public:
  LinkResolver(std::string linker, std::vector<std::string> flags)
      : linker_(linker), flags_(flags) {}

  ::btool::core::VoidErr Resolve(const ::btool::node::Node &n) override {
    ::btool::util::Cmd cmd(linker_.c_str());
    cmd.Arg("-o");
    cmd.Arg(n.name());
    for (const auto &d : *n.dependencies()) {
      cmd.Arg(d->name());
    }
    for (const auto &f : flags_) {
      cmd.Arg(f);
    }

    std::ostringstream out;
    std::ostringstream err;
    cmd.Stdout(&out);
    cmd.Stderr(&err);

    int ec = cmd.Run();

    DEBUG("out: %s\n", out.str().c_str());
    DEBUG("err: %s\n", err.str().c_str());

    return (ec == 0
                ? ::btool::core::VoidErr::Success()
                : ::btool::core::VoidErr::Failure("failed to run compiler"));
  }

 private:
  std::string linker_;
  std::vector<std::string> flags_;
};

::btool::node::Node::Resolver *ResolverFactoryImpl::NewCompileC(
    const std::vector<std::string> &include_dirs,
    const std::vector<std::string> &flags) {
  auto cr =
      new CompileResolver(compiler_c_, compiler_c_flags_, include_dirs, flags);
  allocations_.push_back(cr);
  return cr;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewCompileCC(
    const std::vector<std::string> &include_dirs,
    const std::vector<std::string> &flags) {
  auto cr = new CompileResolver(compiler_cc_, compiler_cc_flags_, include_dirs,
                                flags);
  allocations_.push_back(cr);
  return cr;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewArchive() {
  return nullptr;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewLinkC(
    const std::vector<std::string> &flags) {
  auto lr = new LinkResolver(linker_c_, flags);
  allocations_.push_back(lr);
  return lr;
}

::btool::node::Node::Resolver *ResolverFactoryImpl::NewLinkCC(
    const std::vector<std::string> &flags) {
  auto lr = new LinkResolver(linker_cc_, flags);
  allocations_.push_back(lr);
  return lr;
}

};  // namespace btool::app::collector::cc
