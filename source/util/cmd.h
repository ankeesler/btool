#ifndef BTOOL_UTIL_CMD_H_
#define BTOOL_UTIL_CMD_H_

#include <iostream>
#include <ostream>
#include <sstream>
#include <string>
#include <vector>

namespace btool::util {

class Cmd {
 public:
  Cmd(std::string path)
      : path_(path), dir_(""), out_(&std::cout), err_(&std::cerr) {}

  void Arg(std::string arg) { args_.push_back(arg); }
  void Dir(std::string dir) { dir_ = dir; }
  const std::string &Dir() const { return dir_; }

  void Stdout(std::ostream *out) { out_ = out; }
  void Stderr(std::ostream *err) { err_ = err; }

  int Run(void);

  friend std::ostream &operator<<(std::ostream &os, const Cmd &cmd);

 private:
  int RunChild(int stdout_fds[2], int stderr_fds[2]);
  int RunParent(int child_pid, int child_stdout_fds[2],
                int child_stderr_fds[2]);

  std::string path_;
  std::string dir_;
  std::vector<std::string> args_;
  std::ostream *out_, *err_;
};

std::ostream &operator<<(std::ostream &os, const Cmd &cmd);

};  // namespace btool::util

#endif  // BTOOL_UTIL_CMD_H_
