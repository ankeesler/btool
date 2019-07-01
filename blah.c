#include "blah.h"

#include <stdlib.h>
#include <string.h>

blah_t *blah_new(const char *path) {
  blah_t *blah = (blah_t *)malloc(sizeof(blah_t));
  blah->path = path;
  blah->dependencies = (struct blah_list_t *)blah_list_new(5);
  blah->resolver_f = NULL;
  blah->resolver_ctx = NULL;
  return blah;
}

blah_list_t *blah_list_new(int cap) {
  blah_list_t *l = (blah_list_t *)malloc(sizeof(blah_list_t));
  l->cap = cap;
  l->len = 0;
  l->blahs = (blah_t **)malloc(sizeof(blah_t *) * l->cap);
  return l;
}

void blah_list_add(blah_list_t *l, blah_t *b) {
  if (l->len == l->cap) {
    blah_t **blahs = l->blahs;
    l->cap *= 2;
    l->blahs = (blah_t **)malloc(sizeof(blah_t *) * l->cap);
    memcpy(l->blahs, blahs, sizeof(blah_t *) * l->len);
    free(blahs);
  }

  l->blahs[l->len++] = b;
}

blah_t *blah_list_find(blah_list_t *l, const char *path) {
  for (int i = 0; i < l->len; i++) {
    if (strcmp(l->blahs[i]->path, path) == 0) {
      return l->blahs[i];
    }
  }
  return NULL;
}
