.PHONY: default
default: test

%.c.o: %.c %.h
	clang --std=c11 -g -Wall -Werror -O0 -o $@ -c $<

%.c.o: %.c
	clang --std=c11 -g -Wall -Werror -O0 -o $@ -c $<

blah_test: blah.c.o blah_test.c.o log.c.o
	clang -o $@ $^ -lmcgoo

collect_test: collect.c.o collect_test.c.o blah.c.o log.c.o
	clang -o $@ $^ -lmcgoo

btool: main.c.o log.c.o blah.c.o collect.c.o
	clang -o $@ $^

.PHONY: lint
lint:
	clang-format -i $(shell find . -type d -name fixture -prune -o -type f -name "*.[ch]" -print)
	git diff-index --quiet HEAD -- $(find find . -type d -name fixture -prune -o -type f -name "*.[ch]" -print)

.PHONY: clean
clean:
	find . -name "*.o" | xargs rm -f
	rm -f btool

.PHONY: test
test: lint blah_test collect_test btool
	./blah_test
	./collect_test
	./btool --root fixture/BasicC build main.c && ./main

# TODO: vagrant
