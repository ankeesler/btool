#ifndef BTOOL_NODE_PROPERTYSTORE_H_
#define BTOOL_NODE_PROPERTYSTORE_H_

#include <functional>
#include <map>
#include <string>
#include <vector>

namespace btool::node {

// PropertyStore
//
// PropertyStore is a generic map from a string to value. The value can be of
// the following types:
//   std::vector<std::string>
//   std::string
//   bool
class PropertyStore {
 public:
  PropertyStore() = default;

  PropertyStore(const PropertyStore &ps) = default;
  PropertyStore &operator=(const PropertyStore &ps) = default;

  bool operator==(const PropertyStore &ps) const {
    return ps.bool_store_ == bool_store_ && ps.string_store_ == string_store_ &&
           ps.strings_store_ == strings_store_;
  };

  void Write(const std::string &name, bool b) { bool_store_[name] = b; }

  void Write(const std::string &name, std::string s) {
    string_store_[name] = s;
  }

  void Write(const std::string &name, const char *s) {
    string_store_[name] = s;
  }

  void Append(const std::string &name, const std::string &value) {
    strings_store_[name].push_back(value);
  }

  void ForEach(const std::string &name, std::function<void(std::string *)> f) {
    auto it = strings_store_.find(name);
    if (it == strings_store_.end()) {
      return;
    }

    for (std::size_t i = 0; i < it->second.size(); ++i) {
      f(it->second.data() + i);
    }
  }

  void Read(const std::string &name, const bool **value) const {
    auto it = bool_store_.find(name);
    if (it == bool_store_.end()) {
      *value = nullptr;
    } else {
      *value = &it->second;
    }
  }

  void Read(const std::string &name, const std::string **value) const {
    auto it = string_store_.find(name);
    if (it == string_store_.end()) {
      *value = nullptr;
    } else {
      *value = &it->second;
    }
  }

  void Read(const std::string &name,
            const std::vector<std::string> **value) const {
    auto it = strings_store_.find(name);
    if (it == strings_store_.end()) {
      *value = nullptr;
    } else {
      *value = &it->second;
    }
  }

 private:
  std::map<std::string, bool> bool_store_;
  std::map<std::string, std::string> string_store_;
  std::map<std::string, std::vector<std::string>> strings_store_;
};

};  // namespace btool::node

#endif  // BTOOL_NODE_PROPERTYSTORE_H_
