#ifndef BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORYIMPL_H_
#define BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORYIMPL_H_

#include <memory>
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
                      std::vector<std::string> compiler_cc_flags,
                      std::vector<std::string> linker_c_flags,
                      std::vector<std::string> linker_cc_flags)
      : compiler_c_(compiler_c),
        compiler_cc_(compiler_cc),
        archiver_(archiver),
        linker_c_(linker_c),
        linker_cc_(linker_cc),
        compiler_c_flags_(compiler_c_flags),
        compiler_cc_flags_(compiler_cc_flags),
        linker_c_flags_(linker_c_flags),
        linker_cc_flags_(linker_cc_flags) {}

  ::btool::node::Node::Resolver *NewCompileC() override;
  ::btool::node::Node::Resolver *NewCompileCC() override;
  ::btool::node::Node::Resolver *NewArchive() override;
  ::btool::node::Node::Resolver *NewLinkC() override;
  ::btool::node::Node::Resolver *NewLinkCC() override;

 private:
  std::string compiler_c_;
  std::string compiler_cc_;
  std::string archiver_;
  std::string linker_c_;
  std::string linker_cc_;
  std::vector<std::string> compiler_c_flags_;
  std::vector<std::string> compiler_cc_flags_;
  std::vector<std::string> linker_c_flags_;
  std::vector<std::string> linker_cc_flags_;

  std::vector<std::unique_ptr<::btool::node::Node::Resolver>> resolvers_;
};

};  // namespace btool::app::collector::cc

#endif  // BTOOL_APP_COLLECTOR_CC_RESOLVERFACTORYIMPL_H_
