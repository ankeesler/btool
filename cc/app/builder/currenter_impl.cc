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

static bool GetModTime(const ::btool::node::Node &node, ssize_t *ret_mod_time,
                       std::string *ret_err);
static ssize_t ComputeModTimeNS(struct ::timespec *ts);

bool CurrenterImpl::Current(const ::btool::node::Node &node, bool *ret_current,
                            std::string *ret_err) {
  ssize_t node_mod_time;
  std::string err;
  if (!GetModTime(node, &node_mod_time, &err)) {
    *ret_err = ::btool::WrapErr(err, "get mod time");
    return false;
  } else if (node_mod_time == -1) {
    *ret_current = false;
    return true;
  }
  DEBUGS() << "node " << node.name() << " mod time = " << node_mod_time
           << std::endl;

  bool success = true;
  ssize_t latest_mod_time = node_mod_time;
  const ::btool::node::Node *latest_mod_time_node = &node;
  node.Visit([&](const ::btool::node::Node *dep) {
    if (!success) {
      return;
    }

    ssize_t dep_mod_time;
    if (!GetModTime(*dep, &dep_mod_time, &err)) {
      success = false;
      *ret_err = ::btool::WrapErr(err, "get mod time");
      return;
    }

    DEBUGS() << "dep " << dep->name() << " mod time = " << dep_mod_time
             << std::endl;
    if (dep_mod_time > latest_mod_time) {
      latest_mod_time = dep_mod_time;
      latest_mod_time_node = dep;
    }
  });

  *ret_current = latest_mod_time <= node_mod_time;
  if (*ret_current) {
    DEBUGS() << "dep " << latest_mod_time_node->name() << " is newer than node "
             << node.name() << std::endl;
  }

  return success;
}

#ifdef __linux__
#define modtime st_mtim
#elif __APPLE__
#define modtime st_mtimespec
#else
#error "unknown platform"
#endif

static bool GetModTime(const ::btool::node::Node &node, ssize_t *ret_mod_time,
                       std::string *ret_err) {
  struct ::stat s;
  if (::lstat(node.name().c_str(), &s) == -1) {
    if (errno == ENOENT) {
      *ret_mod_time = -1;
      return true;
    } else {
      *ret_err = ::btool::WrapErr(::strerror(errno), "lstat");
      return false;
    }
  }
  *ret_mod_time = ComputeModTimeNS(&s.modtime);
  return true;
}

static ssize_t ComputeModTimeNS(struct ::timespec *ts) {
  return ((ts->tv_sec * 1e9) + ts->tv_nsec);
}

};  // namespace btool::app::builder
