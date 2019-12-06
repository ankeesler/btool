#ifndef BTOOL_UTIL_FS_FS_H_
#define BTOOL_UTIL_FS_FS_H_

#include <errno.h>
#include <sys/stat.h>
#include <cstring>

#include <chrono>
#include <functional>
#include <string>

#include "err.h"

namespace btool::util::fs {

std::string Base(const std::string &path);
std::string Dir(const std::string &path);
std::string Join(const std::string &one, const std::string &two);
std::string Ext(const std::string &path);

std::string TempDir();

std::string ReadFile(const std::string &path);
void WriteFile(const std::string &path, const std::string &content);

void RemoveAll(const std::string &path);

void Mkdir(const std::string &path);
void MkdirAll(const std::string &path);

bool Exists(const std::string &path);
bool IsDir(const std::string &path);

#ifdef __linux__
#define modtime st_mtim
#elif __APPLE__
#define modtime st_mtimespec
#else
#error "unknown platform"
#endif

template <typename Clock, typename Duration>
std::chrono::time_point<Clock, Duration> ModTime(const std::string &path) {
  struct ::stat s;
  if (::stat(path.c_str(), &s) == -1) {
    THROW_ERR("stat: " + std::string(::strerror(errno)));
  }
  Duration d = std::chrono::duration_cast<Duration>(
      std::chrono::seconds(s.modtime.tv_sec) +
      std::chrono::nanoseconds(s.modtime.tv_nsec));
  return std::chrono::time_point<Clock, Duration>(d);
}

// Walk
//
// Walk performs a depth-first walk on the filesystem at the provided root.
//
// Walk will throw any error it encounters when reading filesystem nodes, or
// the error that is thrown from the provided handler function, f.
//
// Walk must be provided a directory as a root!
void Walk(const std::string &root, std::function<void(const std::string &)> f);

};  // namespace btool::util::fs

#endif  // BTOOL_UTIL_FS_FS_H_
