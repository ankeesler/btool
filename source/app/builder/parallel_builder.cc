#include "app/builder/parallel_builder.h"

#include <map>
#include <set>
#include <vector>

#include "err.h"
#include "log.h"
#include "node/node.h"
#include "util/util.h"

namespace btool::app::builder {

static void BuildDepsCounts(
    const ::btool::node::Node &node,
    std::map<const ::btool::node::Node *, std::set<const ::btool::node::Node *>>
        *deps_counts);
static void CollectLeaves(
    std::map<const ::btool::node::Node *, std::set<const ::btool::node::Node *>>
        *deps_counts,
    std::vector<const ::btool::node::Node *> *leaves);
static void UpdateDepsCounts(
    const ::btool::node::Node *node,
    std::map<const ::btool::node::Node *, std::set<const ::btool::node::Node *>>
        *deps_counts);

void ParallelBuilder::Build(const ::btool::node::Node &node) {
  std::map<const ::btool::node::Node *, std::set<const ::btool::node::Node *>>
      deps_counts;
  BuildDepsCounts(node, &deps_counts);

  std::vector<const ::btool::node::Node *> leaves;
  do {
    CollectLeaves(&deps_counts, &leaves);
    for (auto l : leaves) {
      wp_->Submit([this, l]() -> const ::btool::node::Node * {
        return ReallyBuild(l);
      });
    }
    leaves.clear();

    ::btool::Err err;
    const ::btool::node::Node *n = wp_->Receive(&err);
    if (n == nullptr) {
      THROW_ERR("work failed: " + std::string(err.what()));
    }

    UpdateDepsCounts(n, &deps_counts);
  } while (!deps_counts.empty());
}

const ::btool::node::Node *ParallelBuilder::ReallyBuild(
    const ::btool::node::Node *n) {
  bool current;
  auto d = ::btool::util::Time(
      [this, &current, n]() { current = cu_->Current(*n); });
  (void)d;
  // INFOS() << n->name()
  //        << ": current: " << ::btool::util::CommaSeparatedNumber(d.count())
  //        << " ms" << std::endl;

  DEBUG("builder visiting %s, current: %s, resolver = %s\n", n->name().c_str(),
        (current ? "true" : "false"),
        (n->resolver() == nullptr ? "null" : "something"));
  if (n->resolver() != nullptr) {
    ca_->OnPreResolve(*n, current);
    if (!current) {
      n->resolver()->Resolve(*n);
    }
    ca_->OnPostResolve(*n, current);
  }

  return n;
}

static void BuildDepsCounts(
    const ::btool::node::Node &node,
    std::map<const ::btool::node::Node *, std::set<const ::btool::node::Node *>>
        *deps_counts) {
  node.Visit([deps_counts](const ::btool::node::Node *node) {
    auto deps = node->dependencies();
    deps_counts->emplace(node, std::set<const ::btool::node::Node *>(
                                   deps->begin(), deps->end()));
  });
}

static void CollectLeaves(
    std::map<const ::btool::node::Node *, std::set<const ::btool::node::Node *>>
        *deps_counts,
    std::vector<const ::btool::node::Node *> *leaves) {
  for (const auto &it : *deps_counts) {
    if (it.second.empty()) {
      leaves->push_back(it.first);
    }
  }

  for (auto n : *leaves) {
    deps_counts->erase(n);
  }
}

static void UpdateDepsCounts(
    const ::btool::node::Node *node,
    std::map<const ::btool::node::Node *, std::set<const ::btool::node::Node *>>
        *deps_counts) {
  for (auto &it : *deps_counts) {
    it.second.erase(node);
  }
}

};  // namespace btool::app::builder
