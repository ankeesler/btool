#ifndef __PATH_H__
#define __PATH_H__

#include <string.h>

// Returns string in place.
const char *path_base(const char *path);
const char *path_ext(const char *path);

// Allocates new string.
const char *path_new_ext(const char *path, const char *ext);

#define path_is_c(p) (strncmp((p), "c", 1) == 0)

#endif // __PATH_H__
