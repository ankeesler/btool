#include "blah.h"

#include <stdlib.h>
#include <string.h>

#include "log.h"

blah_t *blah_new(const char *path) {
  blah_t *b = (blah_t *)malloc(sizeof(blah_t));

  b->path = path;
  b->dependencies = (struct blah_list_t *)malloc(sizeof(blah_list_t));
  b->resolver_f = NULL;
  b->resolver_ctx = NULL;

  b->next = NULL;

  return b;
}

void blah_list_add(blah_list_t *l, blah_t *b) {
  if (l->head == NULL) {
    l->head = b;
  } else if (strcmp(l->head->path, b->path) > 0) {
    b->next = l->head;
    l->head = b;
  } else {
    blah_t *tmp = l->head;
    while (tmp->next != NULL && strcmp(tmp->next->path, b->path) <= 0) {
      tmp = tmp->next;
    }

    if (tmp->next == NULL) {
      tmp->next = b;
    } else {
      b->next = tmp->next;
      tmp->next = b;
    }
  }
}

blah_t *blah_list_find(blah_list_t *l, const char *path) {
  blah_list_for_each(l, b) {
    if (strcmp(b->path, path) == 0) {
      return b;
    }
  }
  return NULL;
}

void blah_list_log(blah_list_t *l) {
  blah_list_for_each(l, b) { log_printf("%s", b->path); }
}
