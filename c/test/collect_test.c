#include <unit-test.h>

#include "blah.h"
#include "collect.h"
#include "error.h"

static int happy_test(void) {
  blah_list_t l;
  blah_list_init(&l);

  error_t error = collect_blahs("fixture/basic_c/main.c", &l);
  note(error);
  expect(error == NULL);

  blah_t *mainc = blah_list_find(&l, "fixture/basic_c/main.c");
  expect(mainc != NULL);
  blah_t *masterh = blah_list_find((blah_list_t *)mainc->dependencies,
                                   "fixture/basic_c/master.h");
  expect(masterh != NULL);
  // TODO: dependencies...
  expect(blah_list_find((blah_list_t *)mainc->dependencies, "master.h") !=
         NULL);
  expect(blah_list_find((blah_list_t *)mainc->dependencies, "dep_0/dep_0.h") !=
         NULL);
  expect(blah_list_find((blah_list_t *)mainc->dependencies, "dep_1/dep_1.h") !=
         NULL);

  blah_t *maino = blah_list_find(&l, "fixture/basic_c/main.o");
  expect(maino != NULL);
  expect(blah_list_find((blah_list_t *)mainc->dependencies, "main.c") != NULL);

  // TODO: all other files...

  expect(blah_list_find(&l, "dep-0/dep-0.c") != NULL);
  expect(blah_list_find(&l, "dep-0/dep-0.h") != NULL);
  expect(blah_list_find(&l, "dep-0/dep-0.o") != NULL);
  expect(blah_list_find(&l, "dep-1/dep-1.c") != NULL);
  expect(blah_list_find(&l, "dep-1/dep-1.h") != NULL);
  expect(blah_list_find(&l, "dep-1/dep-1.o") != NULL);

  return 0;
}

static int sad_test(void) {
  blah_list_t l;
  error_t error = collect_blahs("fixture/BasicC/maine.c", &l);
  expect(error != NULL);
  return 0;
}

int main(int argc, char *argv[]) {
  announce();
  run(happy_test);
  run(sad_test);
  return 0;
}
