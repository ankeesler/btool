default: test

GOOGLETEST_DIR=/tmp/googletest/googletest-release-1.8.1

clean:
	rm -rf build

build:
	mkdir build

build/main.o: source/main.cc build
	clang++ -o $@ $< -O0 -g -Wall -Werror -c --std=c++11

build/%.o: source/%.cc source/%.h build
	clang++ -o $@ $< -O0 -g -Wall -Werror -c --std=c++11

build/%.o: test/%.cc build
	clang++ -Isource -I$(GOOGLETEST_DIR)/googletest/include -o $@ $< -O0 -g -Wall -Werror -c --std=c++11

build/%.o: source/cli/%.cc source/cli/%.h build
	clang++ -o $@ $< -O0 -g -Wall -Werror -c --std=c++11

build/%.o: test/cli/%.cc build
	clang++ -o $@ $< -O0 -g -Wall -Werror -c --std=c++11

build/btool: build/main.o build/log.o
	clang++ -o $@ $^

build/error_test: build/error_test.o build/error.o
	clang++ -o $@ -lgtest_main -lgtest $^

.PHONY: run
run: build/btool
	./$<

.PHONY: test
test: build/error_test
	build/error_test
