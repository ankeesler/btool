#include "core/cmd.h"

#include <errno.h>
#include <sys/wait.h>
#include <unistd.h>

#include <cstring>
#include <ostream>

#include "core/log.h"

namespace btool::core {

const int kPipeRead = 0;
const int kPipeWrite = 1;

const int kReadBufSizeLog = 10;  // 1KB

static bool Read(std::ostream *os, int fd) {
  while (true) {
    const int buf_size = 1 << kReadBufSizeLog;
    char buf[buf_size];
    ssize_t count = read(fd, buf, buf_size);
    switch (count) {
      case -1:
        DEBUG("read: %s\n", strerror(errno));
        return false;
      case 0:
        return true;
      default:
        DEBUG("read %d bytes\n", count);
        os->write(buf, count);
    }
  }
}

int Cmd::Run(void) {
  int child_stdout[2];
  int child_stderr[2];
  if (pipe(child_stdout) == -1) {
    DEBUG("pipe: %s\n", strerror(errno));
    return -1;
  }
  if (pipe(child_stderr) == -1) {
    DEBUG("pipe: %s\n", strerror(errno));
    return -1;
  }

  int ret = ::fork();
  switch (ret) {
    case -1:
      return -1;
    case 0:
      return RunChild(child_stdout,
                      child_stderr);  // won't actually return ;)
    default:
      return RunParent(ret, child_stdout, child_stderr);
  }
}

int Cmd::RunChild(int stdout_fds[2], int stderr_fds[2]) {
  // stdout
  if (close(stdout_fds[kPipeRead]) == -1) {
    DEBUG("close: %s\n", strerror(errno));
    return -1;
  }
  if (dup2(stdout_fds[kPipeWrite], STDOUT_FILENO) == -1) {
    DEBUG("dup2: %s\n", strerror(errno));
    return -1;
  }
  if (close(stdout_fds[kPipeWrite]) == -1) {
    DEBUG("close: %s\n", strerror(errno));
    return -1;
  }

  // stderr
  if (close(stderr_fds[kPipeRead]) == -1) {
    DEBUG("close: %s\n", strerror(errno));
    return -1;
  }
  if (dup2(stderr_fds[kPipeWrite], STDERR_FILENO) == -1) {
    DEBUG("dup2: %s\n", strerror(errno));
    return -1;
  }
  if (close(stderr_fds[kPipeWrite]) == -1) {
    DEBUG("close: %s\n", strerror(errno));
    return -1;
  }

  std::vector<const char *> args;

  args.push_back(path_);
  for (auto arg : args_) {
    args.push_back(arg);
  }

  // execv wants a NULL at the end of the args array
  args.push_back(NULL);

  // TF is this cast???
  int ret = ::execvp(args[0], (char *const *)args.data());

  // If the above call worked, we won't get here.
  DEBUG("child failed, exiting: %d (%s)\n", ret, strerror(errno));

  ::exit(ret);

  return ret;
}

int Cmd::RunParent(int child_pid, int child_stdout_fds[2],
                   int child_stderr_fds[2]) {
  if (close(child_stdout_fds[kPipeWrite]) == -1) {
    DEBUG("close: %s\n", strerror(errno));
    return -1;
  }
  if (close(child_stderr_fds[kPipeWrite]) == -1) {
    DEBUG("close: %s\n", strerror(errno));
    return -1;
  }

  if (!Read(out_, child_stdout_fds[kPipeRead])) {
    DEBUG("failed to read from child stdout fd %d\n",
          child_stdout_fds[kPipeRead]);
    return -1;
  }

  if (!Read(err_, child_stderr_fds[kPipeRead])) {
    DEBUG("failed to read from child stderr fd %d\n",
          child_stderr_fds[kPipeRead]);
    return -1;
  }

  if (close(child_stdout_fds[kPipeRead]) == -1) {
    DEBUG("close: %s\n", strerror(errno));
    return -1;
  }
  if (close(child_stderr_fds[kPipeRead]) == -1) {
    DEBUG("close: %s\n", strerror(errno));
    return -1;
  }

  int stat;
  while (true) {
    int pid = wait(&stat);
    if (pid == -1) {
      if (errno != EINTR) {
        DEBUG("wait: %s\n", strerror(errno));
        return -1;
      }
    } else if (pid != child_pid) {
      DEBUG("child %d exited with %d, was looking for %d\n", pid, stat,
            child_pid);
    } else {
      break;
    }
  }

  return WEXITSTATUS(stat);
}

};  // namespace btool::core
