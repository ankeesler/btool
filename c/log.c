#include "log.h"

#include <stdarg.h>
#include <stdio.h>

void log_printf(const char *format, ...) {
  va_list list;
  va_start(list, format);
  printf("btool: ");
  vprintf(format, list);
  printf("\n");
  fflush(stdout);
  va_end(list);
}
