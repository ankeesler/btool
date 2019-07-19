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

include_parser_test: include_parser.o include_parser_test.o log.o  str_buf.o
	clang -o $@ $^ -lmcgoo

include_parser_with_main.o: include_parser.c include_parser.h
	$(CC) -o $@ $< $(CFLAGS) -c -DINCLUDE_PARSER_WITH_MAIN

include_parser: include_parser_with_main.o log.o str_buf.o
	clang -o $@ $^

parser_test: parser.o parser_test.o log.o
	clang -o $@ $^ -lmcgoo

path_test: src/path.o test/path_test.o src/log.o
	clang -o $@ $^ -lmcgoo

str_buf_test: str_buf.o str_buf_test.o
	clang -o $@ $^ -lmcgoo

btool: btool.o log.o blah.o collect.o

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
test: lint run_blah_test run_collect_test run_include_parser_test run_path_test run_str_buf_test btool
	./btool --root fixture/BasicC build main.c && ./main

# TODO: vagrant
