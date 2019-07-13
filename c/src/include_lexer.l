%option noyywrap
%option yylineno

%{
#define YYSTYPE const char *

#include "y.tab.h"
%}

%%

[ \t]                        { ; }
#include                     { return T_INCLUDE; }
\"                           { return T_QUOTE; }
[a-zA-Z\-_0-9\./]+           { yylval = strdup(yytext); return T_FILE; }
.                            { return T_ANYTHING; }

%%

void _(void) {
  (void)(yyunput);
  (void)(input);
}
