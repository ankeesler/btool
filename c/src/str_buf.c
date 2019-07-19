#include "str_buf.h"

#include <stdlib.h>
#include <string.h>

str_buf *str_buf_new(void) {
  const int size = 32;

  str_buf *sb = (str_buf *)malloc(sizeof(str_buf));
  sb->buf = malloc(sizeof(char) * size);
  sb->count = 0;
  sb->size = size;

  return sb;
}

void str_buf_add(str_buf *sb, char c) {
  if (sb->count == sb->size) {
    int new_size = sb->size * 2;
    char *new_buf = (char *)malloc(sizeof(char) * new_size);
    memcpy(new_buf, sb->buf, sb->size);
    free(sb->buf);
    sb->buf = new_buf;
    sb->size = new_size;
  }

  sb->buf[sb->count++] = c;
}

char *str_buf_str(str_buf *sb) {
  char *str = (char *)malloc(sizeof(str_buf) * (sb->count + 1));
  memcpy(str, sb->buf, sb->count);
  str[sb->count] = '\0';
  return str;
}

void str_buf_free(str_buf *sb) {
  free(sb->buf);
  free(sb);
}
