%{
#include "include_parser.h"

#include <stdio.h>

#include "log.h"

#define YYSTYPE const char *

extern int yylex();
extern FILE* yyin;
extern void yyerror(const char *);

static error_t error;
static const char **include_buf;
static int include_buf_idx;
static int include_buf_size;

static void add_include(const char *);

%}

%token T_NEWLINE
%token T_QUOTE
%token T_INCLUDE
%token T_FILE

%start lines

%%

lines:
     | lines line
     ;

line: T_NEWLINE
    | T_INCLUDE T_QUOTE T_FILE T_QUOTE T_NEWLINE { add_include($3); }
    ;
%%

static void add_include(const char *i) {
  log_printf("include %s\n", i);
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
  error = e;
}
