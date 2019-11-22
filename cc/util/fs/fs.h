#ifndef BTOOL_UTIL_FS_FS_H_
#define BTOOL_UTIL_FS_FS_H_

#include <functional>
#include <string>

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

bool Exists(const std::string &path);
bool IsDir(const std::string &path);

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
