#include <stdio.h>

#include "error.h"

#define YYDEBUG 1

error_t parse_includes(FILE *f, const char **buf, int *buf_size);
