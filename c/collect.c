#include "collect.h"

#include <stdio.h>
#include <string.h>
#include <sys/errno.h>

#include "error.h"
#include "include_parser.h"

static error_t collect_dependencies(FILE *f, blah_list_t *l);

error_t collect_blahs(const char *target, blah_list_t *l) {
  FILE *f = fopen(target, "r");
  if (f == NULL) {
    return strerror(errno);
  }

  blah_t *b = blah_new(target);
  blah_list_add(l, b);

  error_t e = collect_dependencies(f, l);
  if (e != NULL) {
    return e;
  }

  if (fclose(f) != 0) {
    return strerror(errno);
  }

  return NULL;
}

static error_t collect_dependencies(FILE *f, blah_list_t *l) {
  char *includes[32];
  int includes_size = sizeof(includes) / sizeof(includes[0]);
  error_t e = parse_includes(f, includes, &includes_size);
  if (e != NULL) {
    return e;
  }

  for (int i = 0; i < includes_size; i++) {
    char *include = includes[i];
    // create blah and add to dependency list
    // if include.c exists, collect blahs on that path
    // collect dependencies on this include
  }
}
