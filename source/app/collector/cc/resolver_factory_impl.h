#ifndef BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORYIMPL_H_
#define BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORYIMPL_H_

#include <string>
#include <vector>

#include "app/collector/cc/resolver_factory.h"
#include "node/node.h"

namespace btool::app::collector::cc {

class ResolverFactoryImpl : public ResolverFactory {
 public:
  ResolverFactoryImpl(std::string compiler_c, std::string compiler_cc,
                      std::string archiver, std::string linker_c,
                      std::string linker_cc,
                      std::vector<std::string> compiler_c_flags,
                      std::vector<std::string> compiler_cc_flags)
      : compiler_c_(compiler_c),
        compiler_cc_(compiler_cc),
        archiver_(archiver),
        linker_c_(linker_c),
        linker_cc_(linker_cc),
        compiler_c_flags_(compiler_c_flags),
        compiler_cc_flags_(compiler_cc_flags) {}

  ~ResolverFactoryImpl() {
    for (auto a : allocations_) {
      delete a;
    }
  }

  ::btool::node::Node::Resolver *NewCompileC(
      const std::vector<std::string> &include_dirs,
      const std::vector<std::string> &flags) override;
  ::btool::node::Node::Resolver *NewCompileCC(
      const std::vector<std::string> &include_dirs,
      const std::vector<std::string> &flags) override;
  ::btool::node::Node::Resolver *NewArchive() override;
  ::btool::node::Node::Resolver *NewLinkC(
      const std::vector<std::string> &flags) override;
  ::btool::node::Node::Resolver *NewLinkCC(
      const std::vector<std::string> &flags) override;

 private:
  std::string compiler_c_;
  std::string compiler_cc_;
  std::string archiver_;
  std::string linker_c_;
  std::string linker_cc_;
  std::vector<std::string> compiler_c_flags_;
  std::vector<std::string> compiler_cc_flags_;

  std::vector<::btool::node::Node::Resolver *> allocations_;
};

};  // namespace btool::app::collector::cc

#endif  // BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORYIMPL_H_
