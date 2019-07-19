#ifndef __LOG_H__
#define __LOG_H__

#define log_printf(format, ...)                                                \
  __log_printf(__FILE__, __LINE__, format, __VA_ARGS__)
void __log_printf(const char *file, int line, const char *format, ...);

#endif // __LOG_H__
