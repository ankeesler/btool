# What's the plan?

## Feature Requests

- Targets are per-os/arch
- Timing numbers for the resolution of each node
- Single line UI for collecting and resolving nodes
- Ability to select multiple targets at once
- Run a debugger with -debug
- Cache registry data
- Run tests as a single linked executable
- Parallel builds

# Chores

- Fix for-loop usage to use for_each if possible.
- Use valgrind to check memory leaks (or -fsanitize= with clang?)
- Is it cool to pass strings in by value to copy them over in constructors?
- Error handling is terrible.
- How do we throw errors in an OnSet call?
- Profile - where are the bottle necks?
- Logging is annoying to have to do printf style.

# Bugs

### 00: Including header doesn't bring that header's include paths
- This is because the exe consumer is pulling in a .o file that has not been
  updated by the inc consumer yet, so it is missing include paths...

### 01: Can't put comments to the right of headers

### 02: Transient header dependencies' source objects are not brought in to build