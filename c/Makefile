.PHONY: default
default: test

CC=clang
CFLAGS=--std=c11 -g -Wall -Werror -O0 -Isrc -I.

YACC=bison
YFLAGS=-d --yacc
LEX=flex
LFLAGS=

blah_test: src/blah.o test/blah_test.o src/log.o
	clang -o $@ $^ -lmcgoo

collect_test: src/collect.o test/collect_test.o src/blah.o src/log.o src/path.o src/include_parser.o src/include_lexer.o
	clang -o $@ $^ -lmcgoo

path_test: src/path.o test/path_test.o src/log.o
	clang -o $@ $^ -lmcgoo

include_parser_test: src/include_parser.o src/include_lexer.o test/include_parser_test.o src/log.o
	clang -o $@ $^ -lmcgoo

include_parser: src/include_parser.o src/include_lexer.o src/log.o
	clang -o $@ $^ -ly

btool: src/btool.o src/log.o src/blah.o src/collect.o

.PHONY: lint
lint:
	clang-format -i $(shell find . -type d -name fixture -prune -o -type f -name "*.[ch]" -print)
	#git diff-index --quiet HEAD -- $(find find . -type d -name fixture -prune -o -type f -name "*.[ch]" -print)

.PHONY: clean
clean:
	find . -name "*.o" | xargs rm -f
	rm -f btool
	rm -f *_test
	rm -f *.tab.[ch] include_parser.c include_lexer.c include_parser y.output

run_%: %
	./$<

.PHONY: test
test: lint run_blah_test run_collect_test run_include_parser_test run_path_test btool
	./btool --root fixture/basic_c build main.c && ./main

# TODO: vagrant
