#include "app/collector/cc/exe.h"

#include <cassert>

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

void CollectObjects(::btool::app::collector::Store *s, ::btool::node::Node *n,
                    const std::string &ext,
                    std::vector<::btool::node::Node *> *objs);
void CollectLibraries(::btool::app::collector::Store *s, ::btool::node::Node *n,
                      std::vector<::btool::node::Node *> *libs);
void CollectLinkFlags(::btool::node::Node *n, std::vector<std::string> *flags);

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
      DEBUGS() << "cannot find source for exe " << name << std::endl;
      assert(0);
    }
  }

  std::vector<::btool::node::Node *> objs;
  CollectObjects(s, src, ext, &objs);
  for (auto obj : objs) {
    n->dependencies()->push_back(obj);
  }

  std::vector<::btool::node::Node *> libs;
  CollectLibraries(s, n, &libs);
  for (auto lib : libs) {
    n->dependencies()->push_back(lib);
  }

  std::vector<std::string> flags;
  CollectLinkFlags(n, &flags);

  auto r = (ext == ".cc" ? rf_->NewLinkCC(flags) : rf_->NewLinkC(flags));
  n->set_resolver(r);

  Notify(s, n->name());
}

void CollectObjects(::btool::app::collector::Store *s, ::btool::node::Node *n,
                    const std::string &ext,
                    std::vector<::btool::node::Node *> *objs) {
  DEBUGS() << "collect from src " << n->name() << std::endl;

  auto obj_name = ::btool::util::string::Replace(n->name(), ext, ".o");
  auto obj = s->Get(obj_name);
  if (obj == nullptr) {
    DEBUGS() << "cannot find obj for obj name " << obj_name << std::endl;
    assert(0);
  }
  if (::btool::util::Contains(*objs, obj)) {
    return;
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
        CollectObjects(s, src, ext, objs);
      }
    }
  }
}

void CollectLibraries(::btool::app::collector::Store *s, ::btool::node::Node *n,
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
      DEBUGS() << "unknown lib for lib_name " << lib_name << std::endl;
      assert(0);
    }

    libs->push_back(lib);
  }
}

void CollectLinkFlags(::btool::node::Node *n, std::vector<std::string> *flags) {
  ::btool::app::collector::CollectStringsProperties(
      n, flags,
      [](const ::btool::node::PropertyStore *ps)
          -> const std::vector<std::string> * {
        return Properties::LinkFlags(ps);
      });
}

}  // namespace btool::app::collector::cc
