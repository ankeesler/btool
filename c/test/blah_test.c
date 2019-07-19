#include "blah.h"

#include <unit-test.h>

static int test(void) {
  blah_list_t l;
  blah_t bs[7] = {
    {
      .path = "one",
    },
    {
      .path = "two",
    },
    {
      .path = "three",
    },
    {
      .path = "four",
    },
    {
      .path = "five",
    },
    {
      .path = "six",
    },
    {
      .path = "seven",
    },
  };
  for (int i = 0; i < sizeof(bs) / sizeof(bs[0]); i++) {
    blah_list_add(&l, &bs[i]);
  }

  for (int i = 0; i < sizeof(bs) / sizeof(bs[0]); i++) {
    expectEquals(blah_list_find(&l, bs[i].path), &bs[i]);
  }

  expectEquals(blah_list_find(&l, "zero"), NULL);

  blah_t sorted_bs[7] = {
    {
      .path = "five",
    },
    {
      .path = "four",
    },
    {
      .path = "one",
    },
    {
      .path = "seven",
    },
    {
      .path = "six",
    },
    {
      .path = "three",
    },
    {
      .path = "two",
    },
  };
  int i = 0;
  blah_list_for_each(&l, b) { expectString(b->path, sorted_bs[i++].path); }

  return 0;
}

static int dependencies(void) {
  blah_t *b0 = blah_new("/some/path/to/file.c");
  blah_t *b1 = blah_new("/some/path/to/other_file.c");
  blah_t *b2 = blah_new("file.h");
  blah_t *b3 = blah_new("master.h");

  blah_list_add((blah_list_t *)b0->dependencies, b2);
  blah_list_add((blah_list_t *)b0->dependencies, b3);

  blah_list_add((blah_list_t *)b1->dependencies, b3);

  return 0;
}

int main(int argc, char *argv[]) {
  announce();
  run(test);
  run(dependencies);
  return 0;
}
