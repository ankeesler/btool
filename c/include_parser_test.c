#include "include_parser.h"

#include <stdlib.h>
#include <unit-test.h>

static int test(void) {
  FILE *f = fopen("fixture/file.h", "r");
  expect(f != NULL);

  const char *buf[4];
  int buf_size = sizeof(buf) / sizeof(buf[0]);
  error_t e = parse_includes(f, buf, &buf_size);
  expectEquals(e, NULL);
  expectEquals(buf_size, 7);

  expectString(buf[0], "file.h");
  expectString(buf[1], "dash-file.h");
  expectString(buf[2], "double-dash-file.h");
  expectString(buf[3], "underscore_file.h");

  // The buffer is filled with pointers to the heap that were allocated
  // as a part of the parsing.
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
