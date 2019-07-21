#include "blah.h"

#include <stdlib.h>
#include <string.h>

#include "log.h"

static blah_t *new_blah(const char *path);
static error_t log_blah(blah_t *b);
static blah_graph_node_t *new_graph_node(blah_t *b);
static blah_graph_node_t *graph_add_node(blah_graph_node_t *node, blah_t *b);
static blah_graph_node_t *graph_find_node(blah_graph_node_t *node,
                                          const char *path);
static error_t graph_walk_nodes(blah_graph_node_t *n,
                                error_t (*walk_f)(blah_t *));

blah_t *blah_new(char *path) {
  blah_t *b = (blah_t *)malloc(sizeof(blah_t));

  b->path = path;
  b->dependencies = (struct blah_list_t *)malloc(sizeof(blah_list_t));
  b->resolver_f = NULL;
  b->resolver_ctx = NULL;

  b->next = NULL;

  blah_list_init(b->dependencies);

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

blah_t *blah_list_find_dependency(blah_list_t *l, const char *path) {
  return NULL;
}

void blah_list_log(blah_list_t *l) {
  blah_list_for_each(l, b) { log_printf("%s", b->path); }
}

blah_t *blah_graph_add_node(blah_graph_t *g, const char *path) {
  blah_t *b = blah_graph_find_node(g, path);
  if (b != NULL) {
    return b;
  }

  blah_t *b = new_blah(path);
  if (g->root == NULL) {
    g->root = new_graph_node(b);
    return g->root->b;
  } else {
    return graph_add_node(g->root, b)->b;
  }
}

blah_t *blah_graph_add_edge(blah_graph_t *g, const char *from, const char *to) {
  blah_graph_add_node(g, from);
  blah_t *b_to = blah_graph_add_node(g, to);

  blah_graph_node_t *n = graph_find_node(g->root, from);
  if (n->dependencies == NULL) {
  }
  graph_add_node(n->dependencies, b_to);

  return b_to;
}

blah_t *blah_graph_find_node(blah_graph_t *g, const char *path) {
  blah_graph_node_t *n = graph_find_node(g->root, path);
  return (n == NULL ? NULL : n->b);
}

error_t blah_graph_walk_nodes(blah_graph_t *g, error_t (*walk_f)(blah_t *)) {
  return graph_walk_nodes(g->root, walk_f);
}

void blah_graph_log(blah_graph_t *g) { blah_graph_walk_nodes(g, log_blah); }

static blah_t *new_blah(const char *path) {
  blah_t *b = (blah_t *)malloc(sizeof(blah_t));
  b->path = strdup(path);
  b->resolver_f = NULL;
  b->resolver_ctx = NULL;
  return b;
}

static error_t log_blah(blah_t *b) {
  log_printf("%s", b->path);
  return NULL;
}

static blah_graph_node_t *new_graph_node(blah_t *b) {
  blah_graph_node_t *n = (blah_graph_node_t *)malloc(sizeof(blah_graph_node_t));
  n->b = b;
  n->left = n->right = NULL;
  return n;
}

static blah_graph_node_t *graph_add_node(blah_graph_node_t *node, blah_t *b) {
  int cmp = strcmp(node->b->path, b->path);
  if (cmp == 0) {
    return node;
  } else if (cmp < 0) {
    if (node->left == NULL) {
      node->left = new_graph_node(b);
      return node->left;
    } else {
      return graph_add_node(node->left, b);
    }
  } else if (cmp > 0) {
    if (node->right == NULL) {
      node->right = new_graph_node(b);
      return node->right;
    } else {
      return graph_add_node(node->right, b);
    }
  }
  return NULL;
}

static blah_graph_node_t *graph_find_node(blah_graph_node_t *node,
                                          const char *path) {
  if (node == NULL) {
    return NULL;
  }

  int cmp = strcmp(node->b->path, path);
  if (cmp == 0) {
    return node;
  } else if (cmp < 0) {
    return graph_find_node(node->left, path);
  } else if (cmp > 0) {
    return graph_find_node(node->right, path);
  }

  return NULL;
}

static error_t graph_walk_nodes(blah_graph_node_t *n,
                                error_t (*walk_f)(blah_t *)) {
  error_t e;

  if (n == NULL) {
    return NULL;
  }

  e = graph_walk_nodes(n->left, walk_f);
  if (e != NULL) {
    return e;
  }

  e = (*walk_f)(n->b);
  if (e != NULL) {
    return e;
  }

  e = graph_walk_nodes(n->right, walk_f);
  if (e != NULL) {
    return e;
  }

  return NULL;
}
