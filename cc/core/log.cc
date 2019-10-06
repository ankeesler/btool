#include "log.h"

#include <cstdarg>
#include <cstdio>

namespace btool::core {

static void logf(const char *file, int line, const char *area,
                 const char *format, va_list args);

void Infof(const char *file, int line, const char *format, ...) {
  va_list args;

  va_start(args, format);
  logf(file, line, "info", format, args);
  va_end(args);
}

void Debugf(const char *file, int line, const char *format, ...) {
  va_list args;

  va_start(args, format);
  logf(file, line, "debug", format, args);
  va_end(args);
}

void Errorf(const char *file, int line, const char *format, ...) {
  va_list args;

  va_start(args, format);
  logf(file, line, "error", format, args);
  va_end(args);
}

static void logf(const char *file, int line, const char *area,
                 const char *format, va_list args) {
  ::fprintf(stderr, "btool | %s | %s:%d | ", area, file, line);
  ::vfprintf(stderr, format, args);
  ::fflush(stderr);
}

};  // namespace btool::core
