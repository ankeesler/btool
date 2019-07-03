#ifndef __BLAH_H__
#define __BLAH_H__

#include "error.h"

struct blah_list_t;

typedef struct blah_tag {
  const char *path;
  struct blah_list_t *dependencies;
  error_t (*resolver_f)(struct blah_tag *blah, void *ctx);
  void *resolver_ctx;

  struct blah_tag *next;
} blah_t;

blah_t *blah_new(const char *path);

typedef struct {
  blah_t *head;
} blah_list_t;

void blah_list_add(blah_list_t *l, blah_t *b);
blah_t *blah_list_find(blah_list_t *l, const char *path);
#define blah_list_for_each(l, b)                                               \
  for (blah_t * (b) = (l)->head; (b) != NULL; (b) = (b)->next)
void blah_list_log(blah_list_t *l);

#endif // __BLAH_H__
