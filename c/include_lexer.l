%option noyywrap

%{
#include "y.tab.h"
%}

%%

[ \t]                        { ; }
"#include"                   { return PARSER_INCLUDE; }
\"                           { return PARSER_QUOTE; }
"/"                          { return PARSER_SLASH; }
[a-zA-Z0-9\.\-_]+            { return PARSER_FILE; }
\n                           { return PARSER_NEWLINE; }
.*                           { return PARSER_ANYTHING; }

%%

void _(void) {
  (void)(yyunput);
  (void)(input);
}
