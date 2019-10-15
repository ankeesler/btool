#ifndef BTOOL_APP_COLLECTOR_RESOLVERFACTORY_H_
#define BTOOL_APP_COLLECTOR_RESOLVERFACTORY_H_

#include <string>
#include <vector>

#include "node/node.h"

namespace btool::app::collector {

class ResolverFactory {
 public:
  virtual class ::btool::node::Node::Resolver *NewCompileC(
      std::vector<std::string> include_dirs,
      std::vector<std::string> flags) = 0;
  virtual class ::btool::node::Node::Resolver *NewCompileCC(
      std::vector<std::string> include_dirs,
      std::vector<std::string> flags) = 0;
  virtual class ::btool::node::Node::Resolver *NewArchive() = 0;
  virtual class ::btool::node::Node::Resolver *NewLinkC(
      std::vector<std::string> flags) = 0;
  virtual class ::btool::node::Node::Resolver *NewLinkCC(
      std::vector<std::string> flags) = 0;

  virtual class ::btool::node::Node::Resolver *NewDownload(
      std::string url, std::string sha256) = 0;
  virtual class ::btool::node::Node::Resolver *NewUnzip() = 0;
};

};  // namespace btool::app::collector

#endif  // BTOOL_APP_COLLECTOR_RESOLVERFACTORY_H_
