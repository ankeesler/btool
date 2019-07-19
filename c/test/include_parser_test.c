#include "include_parser.h"

#include <unit-test.h>

static int test(void) {
  FILE *f = fopen("fixture/file.c", "r");
  expect(f != NULL);

  char *buf[4];
  int buf_size = sizeof(buf) / sizeof(buf[0]);
  error_t e = parse_includes(f, buf, &buf_size);
  expectEquals(e, NULL);
  expectEquals(buf_size, 4);

  expectString(buf[0], "file.h");
  expectString(buf[1], "another/path/to/file_underscore.h");
  expectString(buf[2], "file_underscore.h");
  expectString(buf[3], "path/to/file_underscore.h");

  free((void *)buf[0]);
  free((void *)buf[1]);
  free((void *)buf[2]);
  free((void *)buf[3]);

  expectEquals(fclose(f), 0);

  return 0;
}

int main(int argc, char *argv[]) {
  announce();
  run(test);
  return 0;
}
