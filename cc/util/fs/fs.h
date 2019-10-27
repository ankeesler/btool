#ifndef BTOOL_UTIL_FS_FS_H_
#define BTOOL_UTIL_FS_FS_H_

#include <functional>
#include <string>

#include "core/err.h"

namespace btool::util::fs {

std::string Base(const std::string &path);
std::string Dir(const std::string &path);
std::string Join(const std::string &one, const std::string &two);
std::string Ext(const std::string &path);

::btool::core::Err<std::string> TempDir();

::btool::core::Err<std::string> ReadFile(const std::string &path);
::btool::core::VoidErr WriteFile(const std::string &path,
                                 const std::string &content);

::btool::core::VoidErr RemoveAll(const std::string &path);

::btool::core::VoidErr Mkdir(const std::string &path);

::btool::core::Err<bool> Exists(const std::string &path);
::btool::core::Err<bool> IsFile(const std::string &path);

// Walk
//
// Walk performs a depth-first walk on the filesystem at the provided root.
//
// Walk will return any error it encounters when reading filesystem nodes, or
// the error that is returned from the provided handler function, f.
::btool::core::VoidErr Walk(
    const std::string &root,
    std::function<::btool::core::VoidErr(const std::string &)> f);

};  // namespace btool::util::fs

#endif  // BTOOL_UTIL_FS_FS_H_
