.PHONY: default
default: test

CC=clang
CFLAGS=--std=c11 -g -Wall -Werror -O0

blah_test: blah.o blah_test.o log.o
	clang -o $@ $^ -lmcgoo

collect_test: collect.o collect_test.o blah.o log.o
	clang -o $@ $^ -lmcgoo

btool: btool.o log.o blah.o collect.o

.PHONY: lint
lint:
	clang-format -i $(shell find . -type d -name fixture -prune -o -type f -name "*.[ch]" -print)
	git diff-index --quiet HEAD -- $(find find . -type d -name fixture -prune -o -type f -name "*.[ch]" -print)

.PHONY: clean
clean:
	find . -name "*.o" | xargs rm -f
	rm -f btool

run_%: %
	./$<

.PHONY: test
test: lint run_blah_test run_collect_test btool
	./btool --root fixture/BasicC build main.c && ./main

# TODO: vagrant
