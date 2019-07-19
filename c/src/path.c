#include "path.h"

#include <stdlib.h>
#include <string.h>

#include "log.h"

const char *path_base(const char *path) {
  char *c = strrchr(path, '/');
  if (c == NULL) {
    return path;
  } else {
    return c + 1;
  }
}

const char *path_ext(const char *path) {
  char *c = strrchr(path, '.');
  if (c == NULL) {
    return "";
  } else {
    return c + 1;
  }
}

char *path_new_ext(const char *path, const char *ext) {
  char *c = strrchr(path, '.');
  int size = c - path + strlen(ext) + 1;
  char *path_cpy = (char *)malloc(sizeof(char) * size);
  strcpy(path_cpy, path);

  if (c != NULL) {
    strcpy(path_cpy + (c - path) + 1, ext);
  }

  return path_cpy;
}

char *path_dir(const char *path) {
  char *path_cpy = strdup(path);
  char *c = strrchr(path_cpy, '/');
  if (c == NULL) {
    return "";
  } else {
    *c = '\0';
    return path_cpy;
  }
}

char *path_join(const char *one, const char *two) {
  char *path = (char *)malloc(strlen(one) + strlen(two) + 1);
  char *finger = stpcpy(path, one);
  *finger++ = '/';
  strcpy(finger, two);
  return path;
}
