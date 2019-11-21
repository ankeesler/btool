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
  - Pains:
    - Too many lines of code to get return value from
    - Hard to see exactly where an error came from...
    - Annoying to have to declare multiple variable names to just hold err's
    - VoidErr is misleading - what is "void" about it?
    - Every time we update one class, we have to update the other.
  - Solutions:
    - Offer Wrap utility to add on context where error came from
    - VoidErr
      - bool w/ err pointer parameter output value
      - std::unique_ptr<std::string>
      - Err class with no copy constructor, but instead Wrap
    - Err
      - bool with err pointer and return parameter output value

### Want

- Profile - where are the bottle necks?
- Use valgrind to check memory leaks (or -fsanitize= with clang?)
- Check that files are linted in ci.

## Bugs

### 00: Including header doesn't bring that header's include paths
- This is because the exe consumer is pulling in a .o file that has not been
  updated by the inc consumer yet, so it is missing include paths...

### 01: Can't put comments to the right of headers

### 02: Transient header dependencies' source objects are not brought in to build