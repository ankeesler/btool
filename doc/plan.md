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

- The output of the integration tests is so hard to get the actual error from

### Want

- Where are we copying that we don't need to? More pass-by-reference?
- It is bad that we are calling out to cc package stuff from registry package stuff
- Profile - where are the bottle necks?
- Use valgrind to check memory leaks (or -fsanitize= with clang?)
- Check that files are linted in ci.
- Can the string manipulation in these fs functions be cleaned up at all?
- Cmd should really be printed to stream instead of calling .String()
- Add a README to the example directory to show how to run btool to build examples.
- Add a -help flag to show btool flags.
- Grep for TODOs!

## Bugs

### 00: Including header doesn't bring that header's include paths
- This is because the exe consumer is pulling in a .o file that has not been
  updated by the inc consumer yet, so it is missing include paths...

### 01: Can't put comments to the right of headers

### 02: Transient header dependencies' source objects are not brought in to build

### 03: The way we find include paths is terribly fraught with peril. :(
- A file including `node/node.h` brought in the `node/node.h` from yaml-cpp.