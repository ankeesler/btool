#ifndef __BLAH_H__
#define __BLAH_H__

#include "error.h"

struct blah_list_t;

typedef struct blah_tag {
  char *path;
  struct blah_list_t *dependencies;

  error_t (*resolver_f)(struct blah_tag *blah, void *ctx);
  void *resolver_ctx;

  struct blah_tag *next;
} blah_t;

blah_t *blah_new(char *path);

typedef struct {
  blah_t *head;
} blah_list_t;

#define blah_list_init(l)                                                      \
  do {                                                                         \
    ((blah_list_t *)(l))->head = NULL;                                         \
  } while (0);
void blah_list_add(blah_list_t *l, blah_t *b);
blah_t *blah_list_find(blah_list_t *l, const char *path);
#define blah_list_for_each(l, b)                                               \
  for (blah_t * (b) = (l)->head; (b) != NULL; (b) = (b)->next)
void blah_list_log(blah_list_t *l);

typedef struct blah_graph_node_tag {
  blah_t *b;
  struct blah_graph_node_tag *dependencies;
  struct blah_graph_node_tag *left, *right;
} blah_graph_node_t;

typedef struct {
  blah_graph_node_t *root;
} blah_graph_t;

#define blah_graph_init(g)                                                     \
  do {                                                                         \
    ((blah_graph_t *)(g))->root = NULL;                                        \
  } while (0);

// Calls strdup on the provided path.
// Idempotent.
blah_t *blah_graph_add_node(blah_graph_t *g, const char *path);

// Ensures both from and to are created.
// Returns the "to" edge.
// Idempotent.
blah_t *blah_graph_add_edge(blah_graph_t *g, const char *from, const char *to);

blah_t *blah_graph_find_node(blah_graph_t *g, const char *path);

// In sorted order by path.
error_t blah_graph_walk_nodes(blah_graph_t *g, error_t (*walk_f)(blah_t *));

void blah_graph_log(blah_graph_t *g);

#endif // __BLAH_H__
