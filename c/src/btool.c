#include <string.h>
#include <sys/errno.h>
#include <unistd.h>

#include "blah.h"
#include "build.h"
#include "collect.h"
#include "error.h"
#include "log.h"

// TODO: string memory allocation is nuts.
// TODO: use strNcmp and cousins.
// TODO: having to specify (struct blah_list_t *) sometimes is annoying.
// TODO: reuse blah_t instances.

static char *target = "main.c";
static char *root = "fixture/basic_c";

static error_t run(int argc, char *argv[]);

int main(int argc, char *argv[]) {
  error_t e = run(argc, argv);
  if (e != NULL) {
    log_printf("failure: %s", e);
    return 1;
  }

  return 0;
}

static error_t run(int argc, char *argv[]) {
  if (chdir(root) != 0) {
    return strerror(errno);
  }

  blah_list_t l;
  blah_list_init(&l);

  error_t e = collect_blahs(target, &l);
  if (e != NULL) {
    return e;
  }

  e = build_blahs(&l);
  if (e != NULL) {
    return e;
  }

  return NULL;
}
