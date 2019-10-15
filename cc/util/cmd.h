#ifndef BTOOL_UTIL_CMD_H_
#define BTOOL_UTIL_CMD_H_

#include <iostream>
#include <ostream>
#include <string>
#include <vector>

namespace btool::util {

class Cmd {
 public:
  Cmd(const char *path) : path_(path), out_(&std::cout), err_(&std::cerr) {}

  void Arg(const char *arg) { args_.push_back(arg); }

  void Stdout(std::ostream *out) { out_ = out; }
  void Stderr(std::ostream *err) { err_ = err; }

  int Run(void);

 private:
  int RunChild(int stdout_fds[2], int stderr_fds[2]);
  int RunParent(int child_pid, int child_stdout_fds[2],
                int child_stderr_fds[2]);

  const char *path_;
  std::vector<const char *> args_;
  std::ostream *out_, *err_;
};

};  // namespace btool::util

#endif  // BTOOL_UTIL_CMD_H_
