#include "app/collector/properties.h"

#include <functional>
#include <string>
#include <vector>

#include "node/node.h"
#include "node/property_store.h"

namespace btool::app::collector {

const char *Properties::kLocal = "io.btool.collector.local";
const char *Properties::kRoot = "io.btool.collector.root";

const std::vector<std::string> *ReadStringsProperty(
    const ::btool::node::PropertyStore *ps, const std::string &key) {
  const std::vector<std::string> *strings;
  ps->Read(key, &strings);
  return strings;
}

void CollectStringsProperties(const ::btool::node::Node &n,
                              std::vector<std::string> *accumulator,
                              std::function<const std::vector<std::string> *(
                                  const ::btool::node::PropertyStore *)>
                                  f) {
  n.Visit([accumulator, f](const ::btool::node::Node *vn) {
    const std::vector<std::string> *accs = f(vn->property_store());
    if (accs != nullptr) {
      accumulator->insert(accumulator->end(), accs->begin(), accs->end());
    }
  });
}

};  // namespace btool::app::collector
