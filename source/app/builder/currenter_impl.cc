#include "app/builder/currenter_impl.h"

#include <errno.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <sys/syscall.h>
#include <unistd.h>
#include <cstring>

#include <string>

#include "err.h"
#include "log.h"
#include "node/node.h"

namespace btool::app::builder {

static ssize_t GetModTime(const ::btool::node::Node &node);
static ssize_t ComputeModTimeNS(struct ::timespec *ts);

bool CurrenterImpl::Current(const ::btool::node::Node &node) {
  ssize_t node_mod_time = GetModTime(node);
  if (node_mod_time == -1) {
    return false;
  }
  DEBUG("node %s mod time = %ld\n", node.name().c_str(), node_mod_time);

  ssize_t latest_mod_time = 0;
  const ::btool::node::Node *latest_mod_time_node = nullptr;
  node.Visit([&](const ::btool::node::Node *dep) {
    ssize_t dep_mod_time = GetModTime(*dep);
    DEBUG("dep %s mod time = %ld\n", dep->name().c_str(), dep_mod_time);
    if (dep_mod_time > latest_mod_time) {
      latest_mod_time = dep_mod_time;
      latest_mod_time_node = dep;
    }
  });

  if (latest_mod_time > node_mod_time) {
    DEBUG("dep %s is newer than node %s\n",
          latest_mod_time_node->name().c_str(), node.name().c_str());
    return false;
  } else {
    return true;
  }
}

#ifdef __linux__
#define modtime st_mtim
#elif __APPLE__
#define modtime st_mtimespec
#else
#error "unknown platform"
#endif

static ssize_t GetModTime(const ::btool::node::Node &node) {
  struct ::stat s;
  if (::lstat(node.name().c_str(), &s) == -1) {
    if (errno == ENOENT) {
      return -1;
    } else {
      THROW_ERR("couldn't lstat node");
    }
  }

  return ComputeModTimeNS(&s.modtime);
}

static ssize_t ComputeModTimeNS(struct ::timespec *ts) {
  return ((ts->tv_sec * 1e9) + ts->tv_nsec);
}

};  // namespace btool::app::builder
