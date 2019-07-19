#include "str_buf.h"

#include <unit-test.h>

static int t(void) {
  str_buf *sb = str_buf_new();
  char *s = str_buf_str(sb);
  expectString(s, "");
  free(s);

  str_buf_add(sb, 'a');
  str_buf_add(sb, 'b');
  str_buf_add(sb, 'c');
  s = str_buf_str(sb);
  expectString(s, "abc");
  free(s);

  str_buf_add(sb, 'd');
  str_buf_add(sb, 'e');
  str_buf_add(sb, 'f');
  s = str_buf_str(sb);
  expectString(s, "abcdef");
  free(s);

  str_buf_free(sb);

  return 0;
}

int main(int argc, char *argv[]) {
  announce();
  run(t);
  return 0;
}
