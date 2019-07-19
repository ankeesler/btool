#include <unit-test.h>

#include "blah.h"
#include "build.h"

static int resolver_ctx = 5;
static error_t resolve(struct blah_tag *b, void *ctx) {
  int *int_ctx = (int *)ctx;
  expectEquals(*int_ctx, 5);
  *int_ctx = 10;
  return NULL;
}

static int dependency_resolver_ctx = 100;
static error_t dependency_resolve(struct blah_tag *b, void *ctx) {
  int *int_ctx = (int *)ctx;
  expectEquals(*int_ctx, 100);
  *int_ctx = 200;
  return NULL;
}

static int t(void) {
  blah_list_t l;
  blah_list_init(&l);

  blah_t *something = blah_new("something");
  something->resolver_f = resolve;
  something->resolver_ctx = &resolver_ctx;

  blah_t *dependency = blah_new("dependency");
  dependency->resolver_f = dependency_resolve;
  dependency->resolver_ctx = &dependency_resolver_ctx;

  blah_list_add(&l, something);
  blah_list_add((blah_list_t *)something->dependencies, dependency);

  error_t e = build_blahs(&l);
  expect(e == NULL);

  expectEquals(resolver_ctx, 10);
  expectEquals(dependency_resolver_ctx, 200);

  return 0;
}

int main(int argc, char *argv[]) {
  announce();
  run(t);
  return 0;
}
