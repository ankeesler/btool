#ifndef BTOOL_APP_COLLECTOR_REGISTRY_YAML_FILE_CACHE_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_YAML_FILE_CACHE_H_

#include <fstream>
#include <string>

#include "app/collector/registry/registry.h"
#include "app/collector/registry/serializer.h"
#include "util/cache.h"
#include "util/fs/fs.h"

namespace btool::app::collector::registry {

template <typename T>
class YamlFileCache : public ::btool::util::Cache<T> {
 public:
  YamlFileCache(Serializer<T> *s, std::string dir) : s_(s), dir_(dir) {}

  bool Get(const std::string &name, T *t) override {
    std::string file = ::btool::util::fs::Join(dir_, name);
    if (::btool::util::fs::Exists(file)) {
      std::ifstream ifs(file);
      s_->Unmarshal(&ifs, t);
      return true;
    } else {
      return false;
    }
  }

  void Set(const std::string &name, const T &t) override {
    std::string file = ::btool::util::fs::Join(dir_, name);
    ::btool::util::fs::MkdirAll(::btool::util::fs::Dir(file));
    std::ofstream ofs(file);
    s_->Marshal(&ofs, t);
  }

 private:
  Serializer<T> *s_;
  std::string dir_;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_YAML_FILE_CACHE_H_
