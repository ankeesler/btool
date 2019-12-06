#ifndef BTOOL_APP_COLLECTOR_REGISTRY_YAML_FILE_CACHE_H_
#define BTOOL_APP_COLLECTOR_REGISTRY_YAML_FILE_CACHE_H_

#include <chrono>
#include <fstream>
#include <string>

#include "app/collector/registry/registry.h"
#include "app/collector/registry/serializer.h"
#include "log.h"
#include "util/cache.h"
#include "util/fs/fs.h"

namespace btool::app::collector::registry {

template <typename T>
class YamlFileCache : public ::btool::util::Cache<T> {
 public:
  YamlFileCache(Serializer<T> *s, std::string dir, std::chrono::seconds timeout)
      : s_(s), dir_(dir), timeout_(timeout) {}

  bool Get(const std::string &name, T *t) override {
    std::string file = ::btool::util::fs::Join(dir_, name);
    if (::btool::util::fs::Exists(file) && !IsTimedOut(file)) {
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
  bool IsTimedOut(const std::string &file) {
    auto mod_time = ::btool::util::fs::ModTime<std::chrono::system_clock,
                                               std::chrono::nanoseconds>(file);
    std::chrono::system_clock::time_point now =
        std::chrono::system_clock::now();
    std::chrono::seconds delta =
        std::chrono::duration_cast<std::chrono::seconds>(now - mod_time);
    return delta >= timeout_;
  }

  Serializer<T> *s_;
  std::string dir_;
  std::chrono::seconds timeout_;
};

};  // namespace btool::app::collector::registry

#endif  // BTOOL_APP_COLLECTOR_REGISTRY_YAML_FILE_CACHE_H_
