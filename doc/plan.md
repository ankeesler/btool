# What's the plan?

## Features

### Need

- Ability to run tests

### Pain

- Targets are per-os/arch
- Ability to select multiple targets at once
- Run a debugger with -debug
- After switching to tar-ed archives, the source files are never up to date :(

### Want

- Cache registry data
- When I clean, I don't want to delete the stuff in the cache
- Parallel builds
- Single line UI for collecting and resolving nodes
- Timing numbers for the resolution of each node
- Run tests as a single linked executable

## Chores

### Need

### Pain

### Want

- Where are we copying that we don't need to? More pass-by-reference?
- Profile - where are the bottle necks?
- Use valgrind to check memory leaks (or -fsanitize= with clang?)
- Check that files are linted in ci.

## Bugs

### 00: Including header doesn't bring that header's include paths
- This is because the exe consumer is pulling in a .o file that has not been
  updated by the inc consumer yet, so it is missing include paths...

### 01: Can't put comments to the right of headers

### 02: Transient header dependencies' source objects are not brought in to build

### 03: The way we find include paths is terribly fraught with peril. :(
- A file including `node/node.h` brought in the `node/node.h` from yaml-cpp.