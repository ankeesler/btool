#ifndef __INCLUDE_PARSER_H__
#define __INCLUDE_PARSER_H__

#include <stdio.h>

#include "error.h"

error_t parse_includes(FILE *f, char **buf, int *buf_size);

#endif // __INCLUDE_PARSER_H__
