#include "app/collector/cc/exe.h"

#include <string>
#include <vector>

#include "app/collector/cc/properties.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "log.h"
#include "node/node.h"
#include "util/fs/fs.h"
#include "util/string/string.h"
#include "util/util.h"

namespace btool::app::collector::cc {

void Exe::OnNotify(::btool::app::collector::Store *s, const std::string &name) {
  if (::btool::util::fs::Ext(name) != "") {
    return;
  }

  auto n = s->Get(name);
  if (n == nullptr) {
    return;
  }

  std::string ext = ".c";
  auto src = s->Get(name + ext);
  if (src == nullptr) {
    ext = ".cc";
    src = s->Get(name + ext);
    if (src == nullptr) {
      AddError("cannot find source for exe " + name);
      return;
    }
  }

  std::vector<::btool::node::Node *> objs;
  if (!CollectObjects(s, src, ext, &objs)) {
    return;
  }
  for (auto obj : objs) {
    n->dependencies()->push_back(obj);
  }

  std::vector<::btool::node::Node *> libs;
  if (!CollectLibraries(s, *n, &libs)) {
    return;
  }
  for (auto lib : libs) {
    n->dependencies()->push_back(lib);
  }

  auto r = (ext == ".cc" ? rf_->NewLinkCC() : rf_->NewLinkC());
  n->set_resolver(r);

  Notify(s, n->name());
}

bool Exe::CollectObjects(::btool::app::collector::Store *s,
                         ::btool::node::Node *n, const std::string &ext,
                         std::vector<::btool::node::Node *> *objs) {
  DEBUGS() << "collect from src " << n->name() << std::endl;

  auto obj_name = ::btool::util::string::Replace(n->name(), ext, ".o");
  auto obj = s->Get(obj_name);
  if (obj == nullptr) {
    AddError("cannot find obj for obj name " + obj_name);
    return false;
  }
  if (::btool::util::Contains(*objs, obj)) {
    return true;
  }
  objs->push_back(obj);

  for (auto d : *n->dependencies()) {
    DEBUGS() << "considering dependency " << d->name() << std::endl;
    if (::btool::util::string::HasSuffix(d->name().c_str(), ".h")) {
      auto src_name = ::btool::util::string::Replace(d->name(), ".h", ext);
      auto src = s->Get(src_name);
      if (src == nullptr) {
        DEBUGS() << "no src for src_name " << src_name << std::endl;
      } else {
        bool success = CollectObjects(s, src, ext, objs);
        if (!success) {
          return false;
        }
      }
    }
  }

  return true;
}

bool Exe::CollectLibraries(::btool::app::collector::Store *s,
                           const ::btool::node::Node &n,
                           std::vector<::btool::node::Node *> *libs) {
  std::vector<std::string> lib_names;
  ::btool::app::collector::CollectStringsProperties(
      n, &lib_names,
      [](const ::btool::node::PropertyStore *ps)
          -> const std::vector<std::string> * {
        return Properties::Libraries(ps);
      });

  for (const auto &lib_name : lib_names) {
    auto lib = s->Get(lib_name);
    if (lib == nullptr) {
      AddError("unknown lib for lib_name " + lib_name);
      return false;
    }

    libs->push_back(lib);
  }

  return true;
}

}  // namespace btool::app::collector::cc
