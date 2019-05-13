#ifndef BTOOL_FS_IMPL_H_
#define BTOOL_FS_IMPL_H_

#include <string>

#include "error.h"
#include "log.h"
#include "fs.h"

namespace btool {

class FSImpl : public FS {
public:
  FSImpl(const std::string& root): log_(new Log("fs_impl")), root_(root) { }

  Error WriteFile(const std::string& file, const std::string& contents);

private:
  btool::Log *log_;
  std::string root_;
};

}; // namespace btool

#endif // BTOOL_CREATE_CLASS_COMMAND_H_
