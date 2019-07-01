#include "collect.h"

#include <stdio.h>
#include <string.h>
#include <sys/errno.h>

#include "error.h"

error_t collect_blahs(const char *root, const char *target, blah_list_t *l) {
  const char *p = root; // TODO: join me with target

  FILE *f = fopen(p, "r");
  if (f == NULL) {
    return strerror(errno);
  }

  blah_t *b = blah_new(target);
  blah_list_add(l, b);

  if (fclose(f) != 0) {
    return strerror(errno);
  }

  return NULL;
}
