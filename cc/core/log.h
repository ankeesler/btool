#ifndef BTOOL_CORE_LOG_H_
#define BTOOL_CORE_LOG_H_

namespace btool::core {

void Debugf(const char *, int, const char *, ...);

#define DEBUG(f, ...) Debugf(__FILE__, __LINE__, f, __VA_ARGS__)

};  // namespace btool::core

#endif  // BTOOL_CORE_LOG_H_
