%{
#include "include_parser.h"

#include <stdio.h>

#include "log.h"

#define YYSTYPE const char *

extern int yylex();
extern void yyerror(const char *);
extern FILE* yyin;
extern int yylineno;

static error_t error;
static const char **include_buf;
static int include_buf_idx;
static int include_buf_size;

static void add_include(const char *);

%}

%token T_INCLUDE
%token T_FILE
%token T_ANYTHING

%start lines

%%

lines:
     | lines T_INCLUDE '"' T_FILE '"' { add_include($4); }
     | lines T_INCLUDE '<' T_FILE '>'
     | lines T_ANYTHING
     ;

%%

static void add_include(const char *i) {
  log_printf("include %s", i);

  if (include_buf_idx < include_buf_size) {
    include_buf[include_buf_idx++] = i;
  }
}

error_t parse_includes(FILE *f, const char **buf, int *buf_size) {
  yyin = f;

  error = NULL;
  include_buf = buf;
  include_buf_idx = 0;
  include_buf_size = *buf_size;
  yyparse();

  *buf_size = include_buf_idx;

  return error;
}

void yyerror(const char *e) {
  log_printf("%s, line %d", e, yylineno);
  error = e;
}
