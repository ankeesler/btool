#include "collect.h"

#include <stdio.h>
#include <string.h>
#include <sys/errno.h>

#include "error.h"
#include "include_parser.h"
#include "log.h"
#include "path.h"

#define MAX_INCLUDES 32

static error_t collect_includes(FILE *f, blah_list_t *l);

error_t collect_blahs(const char *target, blah_list_t *l) {
  FILE *f = fopen(target, "r");
  if (f == NULL) {
    return strerror(errno);
  }

  blah_t *b = blah_new(target);
  blah_list_add(l, b);

  error_t e = collect_includes(f, l);
  if (e != NULL) {
    return e;
  }

  if (path_is_c(target)) {
    const char *target_o = path_new_ext(target, "o");
    blah_t *b_o = blah_new(target_o);
    blah_list_add(l, b_o);

    log_printf("creating %s object for %s", target_o, target);
  }

  if (fclose(f) != 0) {
    return strerror(errno);
  }

  return NULL;
}

static error_t collect_includes(FILE *f, blah_list_t *l) {
  const char *includes[MAX_INCLUDES];
  int includes_size = sizeof(includes) / sizeof(includes[0]);
  error_t e = parse_includes(f, includes, &includes_size);
  if (e != NULL) {
    return e;
  }

  for (int i = 0; i < includes_size; i++) {
    blah_t *b = blah_new(includes[i]);
    blah_list_add(l, b);
  }

  return NULL;
}
