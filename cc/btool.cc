#include <cstdlib>

#include <string>

#include "core/flags.h"
#include "core/log.h"

int main(int argc, const char *argv[]) {
  ::btool::core::Flags f;

  bool debug = false;
  f.Bool("debug", &debug);

  std::string err;
  bool success = f.Parse(argc, argv, &err);
  if (!success) {
    ERROR("parse flags: %s\n", err.c_str());
    exit(1);
  }
  DEBUG("debug is %d\n", debug);

  return 0;
}
