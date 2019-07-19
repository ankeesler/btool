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

collect_test: src/collect.o test/collect_test.o src/blah.o src/log.o src/path.o src/include_parser.o src/str_buf.o
	clang -o $@ $^ -lmcgoo

include_parser_test: src/include_parser.o test/include_parser_test.o src/log.o src/str_buf.o
	clang -o $@ $^ -lmcgoo

src/include_parser_with_main.o: src/include_parser.c src/include_parser.h
	$(CC) -o $@ $< $(CFLAGS) -c -DINCLUDE_PARSER_WITH_MAIN

include_parser: src/include_parser_with_main.o src/log.o src/str_buf.o
	clang -o $@ $^

path_test: src/path.o test/path_test.o src/log.o
	clang -o $@ $^ -lmcgoo

str_buf_test: src/str_buf.o test/str_buf_test.o
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
	rm -f *_test include_parser

run_%: %
	./$<

.PHONY: test
test: \
    lint \
    run_blah_test \
    run_include_parser_test \
    run_path_test \
    run_str_buf_test \
    run_collect_test \
    btool
	./btool --root fixture/BasicC build main.c && ./main

# TODO: vagrant
