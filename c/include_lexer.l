%option noyywrap

%{
#include "y.tab.h"
%}

%%

[ \t]                        { ; }
\n                           { return T_NEWLINE; }
\"                           { return T_QUOTE; }
#include                     { return T_INCLUDE; }
.*                           { return T_ANYTHING; }

%%

void _(void) {
  (void)(yyunput);
  (void)(input);
}
