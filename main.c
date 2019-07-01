#include "blah.h"
#include "collect.h"
#include "error.h"
#include "log.h"

static char *root = "/Users/ankeesler/workspace/btool";

static char *target = "main.c";

int main(int argc, char *argv[]) {
  log_printf("start");

  blah_list_t *l = blah_list_new(5);
  error_t error = collect_blahs(root, target, l);
  if (error != NULL) {
    log_printf("error: %s", error);
    return 1;
  }

  return 0;
}
