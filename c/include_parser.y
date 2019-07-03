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

%debug

%%

lines: line
     | lines line
     ;

line: PARSER_NEWLINE
    | include PARSER_NEWLINE
    | PARSER_ANYTHING PARSER_NEWLINE
    ;

include: PARSER_INCLUDE PARSER_QUOTE path PARSER_QUOTE { add_include($3); }
       ;
    
path: PARSER_FILE
    | path PARSER_SLASH PARSER_FILE
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