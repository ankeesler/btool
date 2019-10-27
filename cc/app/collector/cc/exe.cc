#include "app/collector/cc/exe.h"

#include <cassert>

#include <string>
#include <vector>

#include "app/collector/cc/properties.h"
#include "app/collector/properties.h"
#include "app/collector/store.h"
#include "core/log.h"
#include "node/node.h"
#include "util/fs/fs.h"
#include "util/string/string.h"

namespace btool::app::collector::cc {

void CollectObjects(::btool::app::collector::Store *s, ::btool::node::Node *n,
                    const std::string &ext,
                    std::vector<::btool::node::Node *> *objs);
void CollectLibraries(::btool::app::collector::Store *s, ::btool::node::Node *n,
                      std::vector<::btool::node::Node *> *libs);
void CollectLinkFlags(::btool::node::Node *n, std::vector<std::string> *flags);
bool Contains(const std::vector<::btool::node::Node *> &objs,
              const ::btool::node::Node *obj);

void Exe::OnSet(::btool::app::collector::Store *s, const std::string &name) {
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
      DEBUG("cannot find source for exe %s\n", name.c_str());
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
}

void CollectObjects(::btool::app::collector::Store *s, ::btool::node::Node *n,
                    const std::string &ext,
                    std::vector<::btool::node::Node *> *objs) {
  DEBUG("collect from src %s\n", n->name().c_str());

  auto obj_name = ::btool::util::string::Replace(n->name(), ext, ".o");
  auto obj = s->Get(obj_name);
  if (obj == nullptr) {
    DEBUG("cannot find obj for obj name %s\n", obj_name.c_str());
    assert(0);
  }
  if (Contains(*objs, obj)) {
    return;
  }
  objs->push_back(obj);

  for (auto d : *n->dependencies()) {
    DEBUG("considering dependency %s\n", d->name().c_str());
    if (::btool::util::string::HasSuffix(d->name().c_str(), ".h")) {
      auto src_name = ::btool::util::string::Replace(d->name(), ".h", ext);
      auto src = s->Get(src_name);
      if (src == nullptr) {
        DEBUG("no src for src_name %s\n", src_name.c_str());
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

  for (auto lib_name : lib_names) {
    auto lib = s->Get(lib_name);
    if (lib == nullptr) {
      DEBUG("unknown lib for lib_name %s\n", lib_name.c_str());
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

bool Contains(const std::vector<::btool::node::Node *> &objs,
              const ::btool::node::Node *obj) {
  for (auto o : objs) {
    if (o == obj) {
      return true;
    }
  }
  return false;
}

};  // namespace btool::app::collector::cc
