#ifndef BTOOL_LOG_H_
#define BTOOL_LOG_H_

#include <ostream>

namespace btool {

void Debugf(const char *, int, const char *, ...);
void Infof(const char *, int, const char *, ...);
void Errorf(const char *, int, const char *, ...);

#define DEBUG(f, ...) ::btool::Debugf(__FILE__, __LINE__, f, __VA_ARGS__)
#define INFO(f, ...) ::btool::Infof(__FILE__, __LINE__, f, __VA_ARGS__)
#define ERROR(f, ...) ::btool::Errorf(__FILE__, __LINE__, f, __VA_ARGS__)

class Log {
 public:
  static std::ostream *Debug;
  static std::ostream *Info;
  static std::ostream *Error;
};

#define LOGS(s, area) \
  (s) << "btool | " << area << " | " << __FILE__ << ':' << __LINE__ << " | "

#define DEBUGS() LOGS(*::btool::Log::Debug, "debug")
#define INFOS() LOGS(*::btool::Log::Info, "info")
#define ERRORS() LOGS(*::btool::Log::Error, "error")

};  // namespace btool

#endif  // BTOOL_LOG_H_
