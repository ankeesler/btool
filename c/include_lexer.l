%option noyywrap

%{
#define YYSTYPE const char *

#include "y.tab.h"
%}

%%

[ \t]                        { ; }
\n                           { return T_NEWLINE; }
\"                           { return T_QUOTE; }
#include                     { return T_INCLUDE; }
[a-zA-Z\-_0-9\./]+           { yylval = strdup(yytext); return T_FILE; }

%%

void _(void) {
  (void)(yyunput);
  (void)(input);
}
