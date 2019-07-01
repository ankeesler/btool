#include "blah.h"

#include <unit-test.h>

static int test(void) {
  blah_list_t *list = blah_list_new(3);
  blah_t blah[7] = {
    { .path = "one", },
    { .path = "two", },
    { .path = "three", },
    { .path = "four", },
    { .path = "five", },
    { .path = "six", },
    { .path = "seven", },
  };
  for (int i = 0; i < sizeof(blah)/sizeof(blah[0]); i++) {
    blah_list_add(list, &blah[i]);
  }

  for (int i = 0; i < sizeof(blah)/sizeof(blah[0]); i++) {
    expect(blah_list_find(list, blah[i].path) == &blah[i]);
  }

  expect(blah_list_find(list, "zero") == NULL);

  return 0;
}

int main(int argc, char *argv[]) {
  announce();
  run(test);
  return 0;
}
