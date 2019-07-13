#include "parser.h"

#include <unit-test.h>

static int test(void) { return 0; }

int main(int argc, char *argv[]) {
  announce();
  run(test);
  return 0;
}
