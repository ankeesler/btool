#include "blah.h"

#include <string.h>
#include <unit-test.h>

#include "path.h"

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

static blah_t *walk_data[10];
static int walk_data_count = 0;
static error_t walk(blah_t *b) {
  walk_data[walk_data_count++] = b;
  return NULL;
}

static int graph(void) {
  blah_graph_t g;
  blah_graph_init(&g);

  for (int i = 0; i < 3; i++) {
    const char buf[] = {'0' + i, '\0'};

    char *ofile = path_join(buf, "file.o");
    char *cfile = path_join(buf, "file.c");
    char *hfile = path_join(buf, "file.h");

    blah_t *o = blah_graph_add_node(&g, ofile);
    expect(o != NULL);
    expectString(o->path, ofile);
    blah_t *c = blah_graph_add_node(&g, cfile);
    expect(c != NULL);
    expectString(c->path, cfile);
    blah_t *h = blah_graph_add_node(&g, hfile);
    expect(h != NULL);
    expectString(h->path, hfile);

    blah_t *o_again = blah_graph_add_node(&g, ofile);
    expectEquals(o, o_again);
    blah_t *c_again = blah_graph_add_node(&g, cfile);
    expectEquals(c, c_again);
    blah_t *h_again = blah_graph_add_node(&g, hfile);
    expectEquals(h, h_again);

    free(ofile);
    free(cfile);
    free(hfile);
  }

  for (int i = 0; i < 3; i++) {
    const char buf[] = {'0' + i, '\0'};

    char *ofile = path_join(buf, "file.o");
    char *cfile = path_join(buf, "file.c");
    char *hfile = path_join(buf, "file.h");

    blah_t *o = blah_graph_find_node(&g, ofile);
    expect(o != NULL);
    expectString(o->path, ofile);
    blah_t *c = blah_graph_find_node(&g, cfile);
    expect(c != NULL);
    expectString(o->path, ofile);
    blah_t *h = blah_graph_find_node(&g, hfile);
    expect(h != NULL);
    expectString(o->path, ofile);

    blah_t *o_again = blah_graph_add_node(&g, ofile);
    expectEquals(o, o_again);
    blah_t *c_again = blah_graph_add_node(&g, cfile);
    expectEquals(c, c_again);
    blah_t *h_again = blah_graph_add_node(&g, hfile);
    expectEquals(h, h_again);

    free(ofile);
    free(cfile);
    free(hfile);
  }

  blah_t *nope = blah_graph_find_node(&g, "file.c");
  expect(nope == NULL);

  for (int i = 0; i < 3; i++) {
    const char buf[] = {'0' + i, '\0'};

    char *ofile = path_join(buf, "file.o");
    char *cfile = path_join(buf, "file.c");
    char *hfile = path_join(buf, "file.h");

    blah_t *o = blah_graph_find_node(&g, ofile);
    expect(o != NULL);
    expectString(o->path, ofile);
    blah_t *c = blah_graph_find_node(&g, cfile);
    expect(c != NULL);
    expectString(o->path, ofile);
    blah_t *h = blah_graph_find_node(&g, hfile);
    expect(h != NULL);
    expectString(o->path, ofile);

    blah_t *c_edge = blah_graph_add_edge(&g, ofile, cfile);
    expectEquals(c_edge, c);
    blah_t *h_edge = blah_graph_add_edge(&g, cfile, hfile);
    expectEquals(h_edge, h);

    free(ofile);
    free(cfile);
    free(hfile);
  }

  blah_graph_add_edge(&g, "2/file.c", "0/file.h");

  blah_graph_log(&g);

  walk_data_count = 0;
  blah_graph_walk_nodes(&g, walk);
  expectEquals(walk_data_count, 9);
  for (int i = 0; i < 3; i++) {
    const char buf[] = {'0' + i, '\0'};

    char *ofile = path_join(buf, "file.o");
    char *cfile = path_join(buf, "file.c");
    char *hfile = path_join(buf, "file.h");

    expectString(walk_data[i]->path, ofile);
    expectString(walk_data[i]->path, cfile);
    expectString(walk_data[i]->path, hfile);

    free(ofile);
    free(cfile);
    free(hfile);
  }

  return 0;
}

int main(int argc, char *argv[]) {
  announce();
  run(test);
  run(dependencies);
  run(graph);
  return 0;
}
