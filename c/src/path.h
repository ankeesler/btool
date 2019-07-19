#ifndef __PATH_H__
#define __PATH_H__

#include <string.h>
#include <unistd.h>

// Returns string in place.
const char *path_base(const char *path);
const char *path_ext(const char *path);

// Allocates new string.
char *path_new_ext(const char *path, const char *ext);
char *path_dir(const char *path);
char *path_join(const char *one, const char *two);

#define path_is_c(p) (strcmp(path_ext((p)), "c") == 0)
#define path_exists(p) (access((p), R_OK) == 0)

#endif // __PATH_H__
