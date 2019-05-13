#include "fs_impl.h"

#include <fstream>
#include <string>

#include "error.h"

namespace btool {

Error FSImpl::WriteFile(const std::string& path, const std::string& contents) {
  const std::string full_path = root_ + "/" + path;
  log_->Debugf("writing contents '%s' to path '%s'\n",
               contents.c_str(),
               full_path.c_str());
  
  std::ofstream file(full_path);
  if (!file) {
    const std::string msg = "could not open file at path " + full_path;
    return Error::Create(msg.c_str());
  }

  if (!(file << contents)) {
    file.close();
    const std::string msg = "could not write to file at path " + full_path;
    return Error::Create(msg.c_str());
  }

  file.close();

  return Error::Success();
}

}; // namespace btool
