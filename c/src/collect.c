#include "collect.h"

#include <stdio.h>
#include <string.h>
#include <sys/errno.h>
#include <unistd.h>

#include "error.h"
#include "include_parser.h"
#include "log.h"
#include "path.h"

#define MAX_INCLUDES 32

#define log_object(o) log_printf("object %s", (o))
#define log_dependency(from, to) log_printf("dependency %s -> %s", (from), (to))

static error_t collect_includes(FILE *f, blah_list_t *l, blah_t *b);
static char *clean_include(char *include, char *dir);

// Given a C file, walk its dependencies.
error_t collect_blahs(const char *target, blah_list_t *l) {
  FILE *f = fopen(target, "r");
  if (f == NULL) {
    return strerror(errno);
  }
  log_printf("collecting from %s", target);

  if (!path_is_c(target)) {
    return "path is not source file";
  }

  const char *target_o = path_new_ext(target, "o");
  if (blah_list_find(l, target_o) != NULL) {
    return NULL;
  }

  blah_t *b_o = blah_new(target_o);
  blah_list_add(l, b_o);
  log_object(target_o);

  blah_t *b_c = blah_new(target);
  blah_list_add((blah_list_t *)b_o->dependencies, b_c);
  log_dependency(target_o, target);

  error_t e = collect_includes(f, l, b_c);
  if (e != NULL) {
    return e;
  }

  if (fclose(f) != 0) {
    return strerror(errno);
  }

  return NULL;
}

static error_t collect_includes(FILE *f, blah_list_t *l, blah_t *b_c) {
  char *includes[MAX_INCLUDES];
  int includes_size = sizeof(includes) / sizeof(includes[0]);
  error_t e = parse_includes(f, includes, &includes_size);
  if (e != NULL) {
    return e;
  }
  log_printf("found %d includes", includes_size);

  for (int i = 0; i < includes_size; i++) {
    char *include = clean_include(includes[i], path_dir(b_c->path));
    if (include == NULL) {
      log_printf("could not clean path %s for source %s", includes[i],
                 b_c->path);
      continue;
    }
    includes[i] = include;

    blah_t *b_h = blah_new(include);
    blah_list_add((blah_list_t *)b_c->dependencies, b_h);
    log_dependency(b_c->path, b_h->path);
  }

  for (int i = 0; i < includes_size; i++) {
    char *source = path_new_ext(includes[i], "c");
    if (path_exists(source) && blah_list_find(l, source) == NULL) {
      e = collect_blahs(source, l);
      if (e != NULL) {
        return e;
      }
    }
  }

  return NULL;
}

static char *clean_include(char *include, char *dir) {
  if (path_exists(include)) {
    return include;
  }

  char *include_with_dir = path_join(dir, include);
  if (path_exists(include_with_dir)) {
    return include_with_dir;
  }

  return NULL;
}
