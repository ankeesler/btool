#include "include_parser.h"

#include <stdio.h>
#include <string.h>

#include "log.h"
#include "str_buf.h"

typedef enum {
  NONE = 0,

  COMMENT_START = 10,
  COMMENT_LINE = 11,
  COMMENT_MULTILINE = 12,
  COMMENT_STOP = 13,

  INCLUDE_START = 20,
  INCLUDE = 21,

} state_e;

static error_t parse_include_line(FILE *f, const char **buf, int *buf_size);
static error_t parse_include(FILE *f, const char **buf, int *buf_size);

error_t parse_includes(FILE *f, const char **buf, int *buf_size) {
  char c = EOF;
  error_t e = NULL;
  state_e prev_s, cur_s = NONE;
  int buf_idx;

  do {
    prev_s = cur_s;
    c = fgetc(f);
    switch (c) {
    case '/':
      switch (prev_s) {
      case NONE:
        cur_s = COMMENT_START;
        break;
      case COMMENT_START:
        cur_s = COMMENT_LINE;
        break;
      case COMMENT_STOP:
        cur_s = NONE;
        break;
      default:;
      }
      break;

    case '*':
      switch (prev_s) {
      case COMMENT_START:
        cur_s = COMMENT_MULTILINE;
        break;
      case COMMENT_MULTILINE:
        cur_s = COMMENT_STOP;
        break;
      default:;
      }
      break;

    case '\n':
      switch (prev_s) {
      case COMMENT_START:
        cur_s = NONE;
        break;
      case COMMENT_LINE:
        cur_s = NONE;
        break;
      default:;
      }
      break;

    case '#':
      switch (prev_s) {
      case NONE:
        e = parse_include_line(f, buf, &buf_idx);
        cur_s = NONE;
        break;
      default:;
      }
      break;
    }

    // log_printf("%c: %d -> %d\n", c, prev_s, cur_s);
  } while (e == NULL && c != EOF && buf_idx < *buf_size);

  *buf_size = buf_idx;

  return NULL;
}

static error_t parse_include_line(FILE *f, const char **buf, int *buf_idx) {
  char include[8]; // "include" + '\0'
  const int include_len = sizeof(include) / sizeof(include[0]) - 1;
  bzero(include, include_len);

  size_t count = fread(include, sizeof(char), include_len, f);
  if (count != include_len) {
    // log_printf("underflow: %d/%d bytes read", count, include_len);
    return NULL;
  }

  if (strcmp(include, "include") != 0) {
    // log_printf("wrong directive: '%s' != 'include'", include);
    return NULL;
  }

  while (1) {
    char c = fgetc(f);
    switch (c) {
    case '<':
    case '\n':
    case EOF:
      // log_printf("saw '%c', bailing", c);
      return NULL;
    case '"':
      return parse_include(f, buf, buf_idx);
    default:
      continue;
    }
  }
}

static error_t parse_include(FILE *f, const char **buf, int *buf_idx) {
  char c;
  str_buf *sb = str_buf_new();
  char *str = NULL;
  error_t e = NULL;

  do {
    c = fgetc(f);
    switch (c) {
    case '"':
      str_buf_add(sb, '\0');
      str = str_buf_str(sb);
      break;
    case '\n':
      e = "include underflow";
      break;
    default:
      str_buf_add(sb, c);
    }
  } while (str == NULL && e == NULL);

  buf[(*buf_idx)++] = str;

  str_buf_free(sb);

  return e;
}

#ifdef INCLUDE_PARSER_WITH_MAIN

int main(int argc, char *argv[]) {
  const char *buf[64];
  int buf_size = sizeof(buf) / sizeof(buf[0]);
  error_t e = parse_includes(stdin, buf, &buf_size);
  if (e != NULL) {
    printf("error: %s\n", e);
    return 1;
  } else {
    for (int i = 0; i < buf_size; i++) {
      printf("include: %s\n", buf[i]);
    }
    return 0;
  }
}

#endif // INCLUDE_PARSER_WITH_MAIN
