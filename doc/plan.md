# What's the plan?

## Features

- Targets are per-os/arch
- When I clean, I don't want to delete the stuff in the cache
- Parallel builds
- Single line UI for collecting and resolving nodes
- Timing numbers for the resolution of each node
- Ability to select multiple targets at once
- Run a debugger with -debug
- Run tests as a single linked executable
- Ensure that downloaded .tar.gz depends on gaggle.yml (i.e., $this)

## Chores

- Resolvers should be namespaced in yml files
- Print out stats of collection, building, running, etc.
- We could use meta-programming to make currenter cache build up-to-date-ness
- Where are we copying that we don't need to? More pass-by-reference?
  - Can the string manipulation in these fs functions be cleaned up at all?
- Profile - where are the bottle necks?
- It is bad that we are calling out to cc package stuff from registry package stuff

## Bugs

### 00: Including header doesn't bring that header's include paths
- This is because the exe consumer is pulling in a .o file that has not been
  updated by the inc consumer yet, so it is missing include paths...

### 02: Transient header dependencies' source objects are not brought in to build

### 03: The way we find include paths is terribly fraught with peril. :(
- A file including `node/node.h` brought in the `node/node.h` from yaml-cpp.