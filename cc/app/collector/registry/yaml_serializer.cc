#include "app/collector/registry/yaml_serializer.h"

#include "yaml-cpp/yaml.h"

#include "err.h"

namespace YAML {

template <>
struct convert<::btool::app::collector::registry::IndexFile> {
  static Node encode(const ::btool::app::collector::registry::IndexFile &f) {
    Node node;
    node["path"] = f.path;
    node["sha256"] = f.sha256;
    return node;
  }

  static bool decode(const Node &node,
                     ::btool::app::collector::registry::IndexFile &f) {
    if (!node.IsMap() || node.size() != 2) {
      return false;
    }

    f.path = node["path"].as<std::string>();
    f.sha256 = node["sha256"].as<std::string>();

    return true;
  }
};

template <>
struct convert<::btool::app::collector::registry::Index> {
  static Node encode(const ::btool::app::collector::registry::Index &i) {
    Node node;
    for (const auto &f : i.files) {
      node.push_back(f);
    }
    return node;
  }

  static bool decode(const Node &node,
                     ::btool::app::collector::registry::Index &i) {
    if (!node.IsSequence()) {
      return false;
    }

    for (const auto &f : node) {
      i.files.push_back(f.as<::btool::app::collector::registry::IndexFile>());
    }

    return true;
  }
};

}  // namespace YAML

namespace btool::app::collector::registry {

void YamlSerializer::UnmarshalIndex(std::istream *is, Index *i) {
  *i = YAML::Load(*is).as<Index>();
}

void YamlSerializer::UnmarshalGaggle(std::istream *is, Gaggle *g) {}

};  // namespace btool::app::collector::registry
