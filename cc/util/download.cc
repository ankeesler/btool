#include "util/download.h"

#include <iostream>
#include <sstream>
#include <string>

#include "err.h"
#include "util/cmd.h"

namespace btool::util {

void Download(const std::string &url, const std::string &file) {
  Cmd cmd("curl");
  cmd.Arg("-L");
  cmd.Arg("-o");
  cmd.Arg(file);
  cmd.Arg(url);

  std::stringstream out;
  std::stringstream err;
  cmd.Stderr(&out);
  cmd.Stderr(&err);

  int ec = cmd.Run();
  if (ec != 0) {
    std::stringstream ss;
    ss << "failed to run curl command: " << cmd.String() << std::endl;
    ss << "exit code: " << ec << std::endl;
    ss << "stdout: " << out.str() << std::endl;
    ss << "stderr: " << err.str() << std::endl;
    THROW_ERR(ss.str());
  }
}

};  // namespace btool::util
