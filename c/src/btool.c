#include "blah.h"
#include "collect.h"
#include "error.h"
#include "log.h"

// TODO: string memory allocation is nuts.
// TODO: use strNcmp and cousins.
// TODO: having to specify (struct blah_list_t *) sometimes is annoying.
// TODO: reuse blah_t instances.

static char *target = "main.c";

int main(int argc, char *argv[]) {
  log_printf("start");

  blah_list_t l;
  blah_list_init(&l);

  error_t error = collect_blahs(target, &l);
  if (error != NULL) {
    log_printf("error: %s", error);
    return 1;
  }

  return 0;
}
