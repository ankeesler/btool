default: test

GOOGLETEST_DIR=/tmp/googletest/googletest-release-1.8.1
CC=g++

clean:
	rm -rf build

build:
	mkdir build

build/main.o: source/main.cc build
	$(CC) -Isource -o $@ $< -O0 -g -Wall -Werror -c --std=c++11

build/%.o: source/%.cc source/%.h build
	$(CC) -Isource -o $@ $< -O0 -g -Wall -Werror -c --std=c++11

build/%.o: test/%.cc build
	$(CC) -Isource -o $@ $< -O0 -g -Wall -Werror -c --std=c++11

build/%.o: source/cli/%.cc source/cli/%.h build
	$(CC) -Isource -o $@ $< -O0 -g -Wall -Werror -c --std=c++11

build/%.o: test/cli/%.cc build
	$(CC) -Isource -o $@ $< -O0 -g -Wall -Werror -c --std=c++11

build/btool: \
		build/main.o \
		build/log.o \
		build/error.o \
		build/cli.o \
		build/create_class_command.o
	$(CC) -o $@ $^

build/error_test: build/error_test.o build/error.o
	$(CC) -o $@ $^ -lgmock_main -lgmock -lgtest

build/cli_test: build/cli_test.o build/cli.o build/error.o build/log.o
	$(CC) -o $@ $^ -lgmock_main -lgmock -lgtest

build/create_class_command_test: build/create_class_command_test.o build/create_class_command.o build/error.o build/log.o
	$(CC) -o $@ $^ -lgmock_main -lgmock -lgtest

build/fs_impl_test: build/fs_impl_test.o build/fs_impl.o build/error.o build/log.o
	$(CC) -o $@ $^ -lgmock_main -lgmock -lgtest

.PHONY: run
run: build/btool
	./$<

.PHONY: test
test: build/error_test build/cli_test build/create_class_command_test build/fs_impl_test
	build/error_test
	build/cli_test
	build/create_class_command_test
	build/fs_impl_test

.PHONY: containertest
containertest:
	clear
	docker run --rm -it -v $(shell pwd):/etc/btool -w /etc/btool ankeesler/btool make test

.PHONY: containerbuild
containerbuild:
	docker build -t ankeesler/btool .
	docker push ankeesler/btool

.PHONY: containerdebug
containerdebug:
	docker run --rm -it -v $(shell pwd):/etc/btool -w /etc/btool ankeesler/btool sh
