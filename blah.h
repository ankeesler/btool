#ifndef __BLAH_H__
#define __BLAH_H__

#include "error.h"

struct blah_list_t;

typedef struct blah_tag {
  const char *path;
  struct blah_list_t *dependencies;
  error_t (*resolver_f)(struct blah_tag *blah, void *ctx);
  void *resolver_ctx;
} blah_t;

blah_t *blah_new(const char *path);

typedef struct {
  int len, cap;
  blah_t **blahs;
} blah_list_t;

blah_list_t *blah_list_new(int cap);
void blah_list_add(blah_list_t *l, blah_t *b);
blah_t *blah_list_find(blah_list_t *l, const char *path);

#endif // __BLAH_H__
