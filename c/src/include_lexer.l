%option noyywrap
%option yylineno
%option debug

%{
#define YYSTYPE const char *

#include "y.tab.h"
%}

%%

[ \t]                        { ; }
["<>]                        { return yytext[0]; }
#include                     { return T_INCLUDE; }
[a-zA-Z\-_0-9\./]+           { yylval = strdup(yytext); return T_FILE; }
.+                           { return T_ANYTHING; }

%%

void _(void) {
  (void)(yyunput);
  (void)(input);
}
