# What's the plan?

## Feature Requests

- Timing numbers for the resolution of each node
- Single line UI for collecting and resolving nodes
- Ability to select multiple targets at once
- Run a debugger with -debug
- Targets are per-os/arch

# Chores

- Use valgrind to check memory leaks (or -fsanitize= with clang?)
- Is it cool to pass strings in by value to copy them over in constructors?
- Error handling is terrible.
- How do we throw errors in an OnSet call?

# Bugs

### 00: Including header doesn't bring that header's include paths
- This is because the exe consumer is pulling in a .o file that has not been
  updated by the inc consumer yet, so it is missing include paths...

### 01: Can't put comments to the right of headers