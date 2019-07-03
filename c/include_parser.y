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

%token PARSER_INCLUDE
%token PARSER_QUOTE
%token PARSER_SLASH
%token PARSER_FILE
%token PARSER_NEWLINE
%token PARSER_ANYTHING

%start lines

%%

lines: line { printf("aaa '%s' '%s'\n", $$, $1); }
     | lines line { printf("bbb '%s' '%s' '%s'\n", $$, $1, $2); }
     ;

line: PARSER_NEWLINE { printf("ccc '%s' '%s'\n", $$, $1); }
    | include PARSER_NEWLINE { printf("ddd '%s' '%s' '%s'\n", $$, $1, $2); }
    | PARSER_ANYTHING PARSER_NEWLINE { printf("eee '%s' '%s' '%s'\n", $$, $1, $2); }
    ;

include: PARSER_INCLUDE PARSER_QUOTE path PARSER_QUOTE { add_include($3); }
       ;
    
path: PARSER_FILE { printf("fff '%s' '%s'\n", $$, $1); }
    | path PARSER_SLASH PARSER_FILE  { printf("ggg '%s' '%s' '%s'\n", $$, $1, $2); }
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