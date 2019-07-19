#include "build.h"

#include <stdlib.h>

#include "blah.h"
#include "log.h"

static error_t build_blah(blah_t *b);

error_t build_blahs(blah_list_t *l) {
  error_t e;

  blah_list_for_each(l, b) {
    log_printf("build %s", b->path);
    e = build_blah(b);
    if (e != NULL) {
      break;
    }
  }

  return e;
}

static error_t build_blah(blah_t *b) {
  error_t e;
  blah_list_for_each((blah_list_t *)b->dependencies, b_d) {
    log_printf("build dependency %s", b_d->path);
    e = build_blah(b_d);
    if (e != NULL) {
      return e;
    }
  }

  log_printf("resolving %s", b->path);
  e = (*b->resolver_f)((struct blah_tag *)b, b->resolver_ctx);
  if (e != NULL) {
    return e;
  }

  return NULL;
}
