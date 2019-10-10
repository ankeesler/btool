#include "app/builder/currenter_impl.h"

#include <errno.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <sys/syscall.h>
#include <unistd.h>
#include <cstring>

#include <string>

#include "core/err.h"
#include "core/log.h"
#include "node/node.h"

namespace btool::app::builder {

static ::btool::core::Err<ssize_t> GetModTime(const ::btool::node::Node &node);
static ssize_t ComputeModTimeNS(struct ::timespec *ts);

::btool::core::Err<bool> CurrenterImpl::Current(
    const ::btool::node::Node &node) {
  // struct ::stat node_s;
  // if (::lstat(node.Name().c_str(), &node_s) == -1) {
  //  if (errno == ENOENT) {
  //    return ::btool::core::Err<bool>::Success(false);
  //  } else {
  //    DEBUG("lstat: %s\n", strerror(errno));
  //    return ::btool::core::Err<bool>::Failure("couldn't lstat node");
  //  }
  //}
  auto node_mod_time_err = GetModTime(node);
  if (node_mod_time_err) {
    return ::btool::core::Err<bool>::Failure(node_mod_time_err.Msg());
  } else if (node_mod_time_err.Ret() == -1) {
    return ::btool::core::Err<bool>::Success(false);
  }

  ::btool::core::Err<bool> err(true);
  ssize_t latest_mod_time_ns = 0;
  const ::btool::node::Node *latest_mod_time_node = nullptr;
  node.Visit([&](const ::btool::node::Node *dep) {
    if (err) {
      return;
    }

    auto dep_mod_time_err = GetModTime(*dep);
    if (dep_mod_time_err) {
      err = ::btool::core::Err<bool>::Failure(dep_mod_time_err.Msg());
    }

    ssize_t mod_time_ns = dep_mod_time_err.Ret();
    if (mod_time_ns > latest_mod_time_ns) {
      latest_mod_time_ns = mod_time_ns;
      latest_mod_time_node = dep;
    }
  });

  if (err) {
    return err;
  } else if (latest_mod_time_ns > node_mod_time_err.Ret()) {
    DEBUG("dep %s is newer than node %s\n",
          latest_mod_time_node->Name().c_str(), node.Name().c_str());
    return ::btool::core::Err<bool>::Success(false);
  } else {
    return ::btool::core::Err<bool>::Success(true);
  }
}

#ifdef __linux__
static ::btool::core::Err<ssize_t> GetModTime(const ::btool::node::Node &node) {
  int fd = ::open(node.Name().c_str(), O_RDONLY);
  if (fd == -1) {
    if (errno == ENOENT) {
      return ::btool::core::Err<ssize_t>::Success(-1);
    } else {
      DEBUG("lstat: %s\n", strerror(errno));
      return ::btool::core::Err<ssize_t>::Failure("couldn't lstat node");
    }
  }

  struct ::stat s;
  // if (::newfstatat(fd, , &s) == -1) {
  //  DEBUG("lstat: %s\n", strerror(errno));
  //  if (::close(fd) == -1) {
  //    DEBUG("close: %s\n", strerror(errno));
  //  }
  //  return ::btool::core::Err<ssize_t>::Failure("couldn't lstat node");
  //}

  // if (::fstatat(AT_FDCWD, file.c_str(), &s, 0) == -1) {
  // if (::syscall(__NR_newfstatat, AT_FDCWD, file.c_str(), &s, 0) == -1) {
  // if (::lstat(file.c_str(), &s) == -1) {
  if (::fstat(fd, &s) == -1) {
    DEBUG("lstat: %s\n", strerror(errno));
    return ::btool::core::Err<ssize_t>::Failure("lstat");
  }

  if (::close(fd) == -1) {
    DEBUG("close: %s\n", strerror(errno));
  }

  return ComputeModTimeNS(&s.st_mtim);
}
#elif __APPLE__
static ::btool::core::Err<ssize_t> GetModTime(const ::btool::node::Node &node) {
  struct ::stat s;
  if (::lstat(node.Name().c_str(), &s) == -1) {
    if (errno == ENOENT) {
      return ::btool::core::Err<ssize_t>::Success(-1);
    } else {
      DEBUG("lstat: %s\n", strerror(errno));
      return ::btool::core::Err<ssize_t>::Failure("couldn't lstat node");
    }
  }

  return ::btool::core::Err<ssize_t>::Success(
      ComputeModTimeNS(&s.st_mtimespec));
}
#else
#error "Unknown platform"
#endif

static ssize_t ComputeModTimeNS(struct ::timespec *ts) {
  return ((ts->tv_sec * 1e9) + ts->tv_nsec);
}
};  // namespace btool::app::builder
