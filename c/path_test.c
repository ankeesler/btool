#include "path.h"

#include <unit-test.h>

static int base_test(void) {
  expectString("file.h", path_base("file.h"));
  expectString("file.h", path_base("path/to/file.h"));
  expectString("file_with_underscore.c", path_base("path/to/file_with_underscore.c"));
  expectString("file_with_underscore.cc", path_base("path/to/file_with_underscore.cc"));
  return 0;
}

static int ext_test(void) {
  expectString("h", path_ext("file.h"));
  expectString("h", path_ext("path/to/file.h"));
  expectString("c", path_ext("path/to/file_with_underscore.c"));
  expectString("cc", path_ext("path/to/file_with_underscore.cc"));
  return 0;
}

static int new_ext_test(void) {
  expectString("file.c", path_new_ext("file.h", "c"));
  expectString("path/to/file.c", path_new_ext("path/to/file.h", "c"));
  expectString("path/to/file_with_underscore.o", path_new_ext("path/to/file_with_underscore.c", "o"));
  expectString("path/to/file_with_underscore.o", path_new_ext("path/to/file_with_underscore.cc", "o"));
  return 0;
}

int main(int argc, char *argv[]) {
  announce();
  run(base_test);
  run(ext_test);
  run(new_ext_test);
  return 0;
}
