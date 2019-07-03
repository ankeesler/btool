#ifndef __PARSER_H__
#define __PARSER_H__

#include <stdio.h>

#include "error.h"

error_t parse_includes(FILE *f, const char **buf, int *buf_size);

#endif // __PARSER_H__
