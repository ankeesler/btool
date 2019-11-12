# What's the plan?

## Features

### Need

- Ability to run tests

### Pain

- Targets are per-os/arch
- Ability to select multiple targets at once
- Run a debugger with -debug

### Want

- Cache registry data
- Parallel builds
- Single line UI for collecting and resolving nodes
- Timing numbers for the resolution of each node
- Run tests as a single linked executable

## Chores

### Need

### Pain

- Error handling is terrible.
- Logging is annoying to have to do printf style.
- How do we throw errors in an OnSet call?

### Want

- Profile - where are the bottle necks?
- Use valgrind to check memory leaks (or -fsanitize= with clang?)

## Bugs

### 00: Including header doesn't bring that header's include paths
- This is because the exe consumer is pulling in a .o file that has not been
  updated by the inc consumer yet, so it is missing include paths...

### 01: Can't put comments to the right of headers

### 02: Transient header dependencies' source objects are not brought in to build