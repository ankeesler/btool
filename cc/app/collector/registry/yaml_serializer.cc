#include "app/collector/registry/yaml_serializer.h"

#include <vector>

#include "yaml-cpp/yaml.h"

#include "err.h"

namespace YAML {

template <>
struct convert<::btool::app::collector::registry::IndexFile> {
  static bool decode(const Node &node,
                     ::btool::app::collector::registry::IndexFile &f) {
    if (!node.IsMap()) {
      return false;
    }

    f.path = node["path"].as<std::string>();
    f.sha256 = node["sha256"].as<std::string>();

    return true;
  }
};

template <>
struct convert<::btool::app::collector::registry::Index> {
  static bool decode(const Node &node,
                     ::btool::app::collector::registry::Index &i) {
    if (!node.IsMap()) {
      return false;
    }

    i.files =
        node["files"]
            .as<std::vector<::btool::app::collector::registry::IndexFile>>();

    return true;
  }
};

template <>
struct convert<::btool::node::PropertyStore> {
  static bool decode(const Node &node, ::btool::node::PropertyStore &ps) {
    if (!node.IsMap()) {
      return false;
    }

    for (const auto &it : node) {
      auto key = it.first.as<std::string>();
      auto value = node[key];
      switch (value.Type()) {
        case NodeType::Scalar: {
          auto s = value.as<std::string>();
          if (s == "true" || s == "false") {
            ps.Write(key, value.as<bool>());
          } else {
            ps.Write(key, value.as<std::string>());
          }
          break;
        }

        case NodeType::Sequence:
          for (const auto &s : value) {
            ps.Append(key, s.as<std::string>());
          }
          break;

        case NodeType::Map:
        case NodeType::Undefined:
        case NodeType::Null:
        default:
          return false;
      }
    }

    return true;
  }
};

template <>
struct convert<::btool::app::collector::registry::Resolver> {
  static bool decode(const Node &node,
                     ::btool::app::collector::registry::Resolver &r) {
    if (!node.IsMap()) {
      return false;
    }

    r.name = node["name"].as<std::string>();
    if (node["config"]) {
      r.config = node["config"].as<::btool::node::PropertyStore>();
    }

    return true;
  }
};

template <>
struct convert<::btool::app::collector::registry::Node> {
  static bool decode(const Node &node,
                     ::btool::app::collector::registry::Node &n) {
    if (!node.IsMap()) {
      return false;
    }

    n.name = node["name"].as<std::string>();
    if (node["dependencies"]) {
      n.dependencies = node["dependencies"].as<std::vector<std::string>>();
    }
    if (node["labels"]) {
      n.labels = node["labels"].as<::btool::node::PropertyStore>();
    }
    if (node["resolver"]) {
      n.resolver =
          node["resolver"].as<::btool::app::collector::registry::Resolver>();
    }

    return true;
  }
};

template <>
struct convert<::btool::app::collector::registry::Gaggle> {
  static bool decode(const Node &node,
                     ::btool::app::collector::registry::Gaggle &g) {
    if (!node.IsMap()) {
      return false;
    }

    g.nodes = node["nodes"]
                  .as<std::vector<::btool::app::collector::registry::Node>>();

    return true;
  }
};

}  // namespace YAML

namespace btool::app::collector::registry {

void YamlSerializer::UnmarshalIndex(std::istream *is, Index *i) {
  try {
    *i = YAML::Load(*is).as<Index>();
  } catch (const YAML::Exception &e) {
    THROW_ERR("could not unmarshal index: " + std::string(e.what()));
  }
}

void YamlSerializer::UnmarshalGaggle(std::istream *is, Gaggle *g) {
  try {
    *g = YAML::Load(*is).as<Gaggle>();
  } catch (const YAML::Exception &e) {
    THROW_ERR("could not unmarshal gaggle: " + std::string(e.what()));
  }
}

};  // namespace btool::app::collector::registry
