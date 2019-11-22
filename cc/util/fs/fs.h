#ifndef BTOOL_UTIL_FS_FS_H_
#define BTOOL_UTIL_FS_FS_H_

#include <functional>
#include <string>

#include "err.h"

namespace btool::util::fs {

std::string Base(const std::string &path);
std::string Dir(const std::string &path);
std::string Join(const std::string &one, const std::string &two);
std::string Ext(const std::string &path);

bool TempDir(std::string *ret_dir, std::string *ret_err);

bool ReadFile(const std::string &path, std::string *ret_content,
              std::string *ret_err);
bool WriteFile(const std::string &path, const std::string &content,
               std::string *ret_err);

bool RemoveAll(const std::string &path, std::string *ret_err);

bool Mkdir(const std::string &path, std::string *ret_err);

bool Exists(const std::string &path, bool *ret_exists, std::string *ret_err);
bool IsDir(const std::string &path, bool *ret_is, std::string *ret_err);

// Walk
//
// Walk performs a depth-first walk on the filesystem at the provided root.
//
// Walk will return any error it encounters when reading filesystem nodes, or
// the error that is returned from the provided handler function, f.
//
// Walk must be provided a directory as a root!
bool Walk(const std::string &root,
          std::function<bool(const std::string &, std::string *)> f,
          std::string *ret_err);

};  // namespace btool::util::fs

#endif  // BTOOL_UTIL_FS_FS_H_
