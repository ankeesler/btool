#include "app/collector/registry/yaml_serializer.h"

#include <istream>
#include <ostream>
#include <vector>

#include "err.h"
#include "yaml-cpp/yaml.h"

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
  static Node encode(const ::btool::app::collector::registry::Index &i) {
    Node node;
    node["files"] = i.files;
    return node;
  }

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
  static Node encode(const ::btool::node::PropertyStore &ps) {
    Node node;
    ps.ForEach([&node, ps](const std::string &name,
                           ::btool::node::PropertyStore::Type type) {
      switch (type) {
        case ::btool::node::PropertyStore::kBool: {
          const bool *b = nullptr;
          ps.Read(name, &b);
          node[name] = *b;
          break;
        }
        case ::btool::node::PropertyStore::kString: {
          const std::string *s = nullptr;
          ps.Read(name, &s);
          node[name] = *s;
          break;
        }
        case ::btool::node::PropertyStore::kStrings: {
          const std::vector<std::string> *ss = nullptr;
          ps.Read(name, &ss);
          node[name] = *ss;
          break;
        }
      }
    });
    return node;
  }

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
  static Node encode(const ::btool::app::collector::registry::Resolver &r) {
    Node node;
    node["name"] = r.name;
    node["config"] = r.config;
    return node;
  }

  static bool decode(const Node &node,
                     ::btool::app::collector::registry::Resolver &r) {
    if (!node.IsMap()) {
      return false;
    }

    r.name = node["name"].as<std::string>();
    if (node["config"] && node["config"].Type() != NodeType::Null) {
      r.config = node["config"].as<::btool::node::PropertyStore>();
    }

    return true;
  }
};

template <>
struct convert<::btool::app::collector::registry::Node> {
  static Node encode(const ::btool::app::collector::registry::Node &n) {
    Node node;
    node["name"] = n.name;
    node["dependencies"] = n.dependencies;
    node["labels"] = n.labels;
    node["resolver"] = n.resolver;
    return node;
  }

  static bool decode(const Node &node,
                     ::btool::app::collector::registry::Node &n) {
    if (!node.IsMap()) {
      return false;
    }

    n.name = node["name"].as<std::string>();
    if (node["dependencies"] && node["dependencies"].Type() != NodeType::Null) {
      n.dependencies = node["dependencies"].as<std::vector<std::string>>();
    }
    if (node["labels"] && node["labels"].Type() != NodeType::Null) {
      n.labels = node["labels"].as<::btool::node::PropertyStore>();
    }
    if (node["resolver"] && node["resolver"].Type() != NodeType::Null) {
      n.resolver =
          node["resolver"].as<::btool::app::collector::registry::Resolver>();
    }

    return true;
  }
};

template <>
struct convert<::btool::app::collector::registry::Gaggle> {
  static Node encode(const ::btool::app::collector::registry::Gaggle &g) {
    Node node;
    node["nodes"] = g.nodes;
    return node;
  }

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

template <>
void YamlSerializer<Index>::Unmarshal(std::istream *is, Index *i) {
  try {
    *i = YAML::Load(*is).as<Index>();
  }  // namespace YAML
  catch (const YAML::Exception &e) {
    THROW_ERR("could not unmarshal index: " + std::string(e.what()));
  }
}

template <>
void YamlSerializer<Index>::Marshal(std::ostream *os, const Index &i) {
  *os << YAML::Node(i);
}

template <>
void YamlSerializer<Gaggle>::Unmarshal(std::istream *is, Gaggle *g) {
  try {
    *g = YAML::Load(*is).as<Gaggle>();
  } catch (const YAML::Exception &e) {
    THROW_ERR("could not unmarshal gaggle: " + std::string(e.what()));
  }
}

template <>
void YamlSerializer<Gaggle>::Marshal(std::ostream *os, const Gaggle &g) {
  *os << YAML::Node(g);
}
};  // namespace btool::app::collector::registry
