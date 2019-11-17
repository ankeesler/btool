#include "log.h"

#include <cstdarg>
#include <cstdio>

#include <fstream>
#include <ostream>

namespace btool {

Log::Level Log::CurrentLevel = Log::kInfo;

std::ostream *Log::Debug = &std::cerr;
std::ostream *Log::Info = &std::cerr;
std::ostream *Log::Error = &std::cerr;

static void logf(const char *file, int line, const char *area,
                 const char *format, va_list args);

void Infof(const char *file, int line, const char *format, ...) {
  if (!Log::IsLevelEnabled(Log::kInfo)) {
    return;
  }

  va_list args;

  va_start(args, format);
  logf(file, line, "info", format, args);
  va_end(args);
}

void Debugf(const char *file, int line, const char *format, ...) {
  if (!Log::IsLevelEnabled(Log::kDebug)) {
    return;
  }

  va_list args;

  va_start(args, format);
  logf(file, line, "debug", format, args);
  va_end(args);
}

void Errorf(const char *file, int line, const char *format, ...) {
  if (!Log::IsLevelEnabled(Log::kError)) {
    return;
  }

  va_list args;

  va_start(args, format);
  logf(file, line, "error", format, args);
  va_end(args);
}

enum Log::Level Log::ParseLevel(const std::string &loglevel) {
  if (loglevel == "debug") {
    return Log::kDebug;
  } else if (loglevel == "info") {
    return Log::kInfo;
  } else if (loglevel == "error") {
    return Log::kError;
  } else {
    return Log::kUnknown;
  }
}

void Log::SetCurrentLevel(Level level) {
  static std::ofstream DevNull("/dev/null");

  Log::CurrentLevel = level;

  Log::Debug = &DevNull;
  Log::Info = &DevNull;
  Log::Error = &DevNull;

  switch (Log::CurrentLevel) {
    case Log::kDebug:
      Log::Debug = &std::cerr;
      // fallthrough
    case Log::kInfo:
      Log::Info = &std::cerr;
      // fallthrough
    case Log::kError:
      Log::Error = &std::cerr;
      // fallthrough
    default:;
  }
}

static void logf(const char *file, int line, const char *area,
                 const char *format, va_list args) {
  ::fprintf(stderr, "btool | %s | %s:%d | ", area, file, line);
  ::vfprintf(stderr, format, args);
  ::fflush(stderr);
}

};  // namespace btool
