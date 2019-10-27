#ifndef BTOOL_NODE_PROPERTYSTORE_H_
#define BTOOL_NODE_PROPERTYSTORE_H_

#include <map>
#include <string>
#include <vector>

namespace btool::node {

// PropertyStore
//
// PropertyStore is a generic map from a string to value. The value can be of
// the following types:
//   std::vector<std::string>
//   bool
class PropertyStore {
 public:
  void Write(const std::string &name, bool b) { bool_store_[name] = b; }

  void Append(const std::string &name, const std::string &value) {
    strings_store_[name].push_back(value);
  }

  void Read(const std::string &name, const bool **value) const {
    auto it = bool_store_.find(name);
    if (it == bool_store_.end()) {
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
  std::map<std::string, std::vector<std::string>> strings_store_;
};

};  // namespace btool::node

#endif  // BTOOL_NODE_PROPERTYSTORE_H_
