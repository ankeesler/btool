%option noyywrap

%{
#include "y.tab.h"
%}

%%

[ \t\n]                      { ; }
#include                     { return PARSER_INCLUDE; }
\"                           { return PARSER_QUOTE; }
.*                           { return PARSER_ANYTHING; }

%%

void _(void) {
  (void)(yyunput);
  (void)(input);
}
