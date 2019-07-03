#include "include_parser.h"

#include <unit-test.h>

static int test(void) {
  FILE *f = fopen("fixture/file.h", "r");
  expect(f != NULL);

  const char *buf[4];
  int buf_size = sizeof(buf) / sizeof(buf[0]);
  error_t e = parse_includes(f, buf, &buf_size);
  expectEquals(e, NULL);
  expectEquals(buf_size, 4);

  expectEquals(fclose(f), 0);

  return 0;
}

int main(int argc, char *argv[]) {
  announce();
  run(test);
  return 0;
}
