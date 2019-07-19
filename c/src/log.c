#include "log.h"

#include <stdarg.h>
#include <stdio.h>

void __log_printf(const char *file, int line, const char *format, ...) {
  va_list list;
  va_start(list, format);
  printf("btool:%s:%d | ", file, line);
  vprintf(format, list);
  printf("\n");
  fflush(stdout);
  va_end(list);
}
