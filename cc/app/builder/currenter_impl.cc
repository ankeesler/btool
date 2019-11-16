#include "app/builder/currenter_impl.h"

#include <errno.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <sys/syscall.h>
#include <unistd.h>
#include <cstring>

#include <string>

#include "err.h"
#include "core/log.h"
#include "node/node.h"

namespace btool::app::builder {

static ::btool::Err<ssize_t> GetModTime(const ::btool::node::Node &node);
static ssize_t ComputeModTimeNS(struct ::timespec *ts);

::btool::Err<bool> CurrenterImpl::Current(
    const ::btool::node::Node &node) {
  auto node_mod_time_err = GetModTime(node);
  if (node_mod_time_err) {
    return ::btool::Err<bool>::Failure(node_mod_time_err.Msg());
  } else if (node_mod_time_err.Ret() == -1) {
    return ::btool::Err<bool>::Success(false);
  }
  DEBUG("node %s mod time = %ld\n", node.name().c_str(),
        node_mod_time_err.Ret());

  ::btool::Err<bool> err(true);
  ssize_t latest_mod_time_ns = 0;
  const ::btool::node::Node *latest_mod_time_node = nullptr;
  node.Visit([&](const ::btool::node::Node *dep) {
    if (err) {
      return;
    }

    auto dep_mod_time_err = GetModTime(*dep);
    if (dep_mod_time_err) {
      err = ::btool::Err<bool>::Failure(dep_mod_time_err.Msg());
    }

    ssize_t mod_time_ns = dep_mod_time_err.Ret();
    DEBUG("dep %s mod time = %ld\n", dep->name().c_str(), mod_time_ns);
    if (mod_time_ns > latest_mod_time_ns) {
      latest_mod_time_ns = mod_time_ns;
      latest_mod_time_node = dep;
    }
  });

  if (err) {
    return err;
  } else if (latest_mod_time_ns > node_mod_time_err.Ret()) {
    DEBUG("dep %s is newer than node %s\n",
          latest_mod_time_node->name().c_str(), node.name().c_str());
    return ::btool::Err<bool>::Success(false);
  } else {
    return ::btool::Err<bool>::Success(true);
  }
}

#ifdef __linux__
#define modtime st_mtim
#elif __APPLE__
#define modtime st_mtimespec
#else
#error "unknown platform"
#endif

static ::btool::Err<ssize_t> GetModTime(const ::btool::node::Node &node) {
  struct ::stat s;
  if (::lstat(node.name().c_str(), &s) == -1) {
    if (errno == ENOENT) {
      return ::btool::Err<ssize_t>::Success(-1);
    } else {
      DEBUG("lstat: %s\n", strerror(errno));
      return ::btool::Err<ssize_t>::Failure("couldn't lstat node");
    }
  }

  return ::btool::Err<ssize_t>::Success(ComputeModTimeNS(&s.modtime));
}

static ssize_t ComputeModTimeNS(struct ::timespec *ts) {
  return ((ts->tv_sec * 1e9) + ts->tv_nsec);
}

};  // namespace btool::app::builder
