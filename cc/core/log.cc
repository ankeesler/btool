#include "log.h"

#include <cstdarg>
#include <cstdio>

namespace btool::core {

void Debugf(const char *file, int line, const char *format, ...) {
  va_list list;

  va_start(list, format);
  fprintf(stderr, "btool | %s:%d | ", file, line);
  vfprintf(stderr, format, list);
  va_end(list);
}

};  // namespace btool::core
