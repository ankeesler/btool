#include <unistd.h>
#include <unit-test.h>

#include "blah.h"
#include "collect.h"
#include "error.h"

static int happy_test(void) {
  blah_list_t l;
  blah_list_init(&l);

  expectEquals(chdir("fixture/basic_c"), 0);

  error_t error = collect_blahs("main.c", &l);
  expect(error == NULL);

  blah_t *maino = blah_list_find(&l, "main.o");
  expect(maino != NULL);
  blah_t *mainc = blah_list_find((blah_list_t *)maino->dependencies, "main.c");
  expect(mainc != NULL);

  blah_t *masterh =
    blah_list_find((blah_list_t *)mainc->dependencies, "master.h");
  expect(masterh != NULL);
  blah_t *dep0h =
    blah_list_find((blah_list_t *)mainc->dependencies, "dep_0/dep_0.h");
  expect(dep0h != NULL);
  blah_t *dep1h =
    blah_list_find((blah_list_t *)mainc->dependencies, "dep_1/dep_1.h");
  expect(dep1h != NULL);

  blah_t *dep0o = blah_list_find(&l, "dep_0/dep_0.o");
  expect(dep0o != NULL);
  blah_t *dep0c =
    blah_list_find((blah_list_t *)dep0o->dependencies, "dep_0/dep_0.c");
  expect(dep0c != NULL);
  dep0h = blah_list_find((blah_list_t *)dep0c->dependencies, "dep_0/dep_0.h");
  expect(dep0h != NULL);

  blah_t *dep1o = blah_list_find(&l, "dep_1/dep_1.o");
  expect(dep1o != NULL);
  blah_t *dep1c =
    blah_list_find((blah_list_t *)dep1o->dependencies, "dep_1/dep_1.c");
  expect(dep1c != NULL);
  dep1h = blah_list_find((blah_list_t *)dep1c->dependencies, "dep_1/dep_1.h");
  expect(dep1h != NULL);
  dep0c = blah_list_find((blah_list_t *)dep1h->dependencies, "dep_0/dep_0.h");
  expect(dep0h != NULL);

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
