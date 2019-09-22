#ifndef BTOOL_CORE_LOG_H_
#define BTOOL_CORE_LOG_H_

namespace btool::core {

void Debugf(const char *, int, const char *, ...);
void Infof(const char *, int, const char *, ...);
void Errorf(const char *, int, const char *, ...);

#define DEBUG(f, ...) ::btool::core::Debugf(__FILE__, __LINE__, f, __VA_ARGS__)
#define INFO(f, ...) ::btool::core::Infof(__FILE__, __LINE__, f, __VA_ARGS__)
#define ERROR(f, ...) ::btool::core::Errorf(__FILE__, __LINE__, f, __VA_ARGS__)

};  // namespace btool::core

#endif  // BTOOL_CORE_LOG_H_
