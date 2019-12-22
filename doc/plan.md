# What's the plan?

## Features

- Targets are per-os/arch
- When I clean, I don't want to delete the stuff in the cache
- Add an -output flag
- Single line UI for collecting and resolving nodes
- Ability to select multiple targets at once
- Run a debugger with -debug (needs Cmd to support stdin)
- Run tests as a single linked executable
- Ensure that downloaded .tar.gz depends on gaggle.yml (i.e., $this)
- People shouldn't have to install transient dependencies themselves: perl, make

## Chores

- We could use meta-programming to make currenter cache build up-to-date-ness
- Where are we copying that we don't need to? More pass-by-reference?
  - Can the string manipulation in these fs functions be cleaned up at all?
- Profile - where are the bottle necks?
- It is bad that we are calling out to cc package stuff from registry package stuff

## Bugs

### 03: The way we find include paths is terribly fraught with peril. :(
- A file including `node/node.h` brought in the `node/node.h` from yaml-cpp.
