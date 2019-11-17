#ifndef BTOOL_LOG_H_
#define BTOOL_LOG_H_

namespace btool {

void Debugf(const char *, int, const char *, ...);
void Infof(const char *, int, const char *, ...);
void Errorf(const char *, int, const char *, ...);

#define DEBUG(f, ...) ::btool::Debugf(__FILE__, __LINE__, f, __VA_ARGS__)
#define INFO(f, ...) ::btool::Infof(__FILE__, __LINE__, f, __VA_ARGS__)
#define ERROR(f, ...) ::btool::Errorf(__FILE__, __LINE__, f, __VA_ARGS__)

};  // namespace btool

#endif  // BTOOL_LOG_H_
