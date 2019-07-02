#include "blah.h"
#include "collect.h"
#include "error.h"
#include "log.h"

static char *target = "main.c";

int main(int argc, char *argv[]) {
  log_printf("start");

  blah_list_t l;
  error_t error = collect_blahs(target, &l);
  if (error != NULL) {
    log_printf("error: %s", error);
    return 1;
  }

  return 0;
}
