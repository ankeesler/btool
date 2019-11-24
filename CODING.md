# Coding Standard

## Go

### Directory Structure

- Each package should only depend on concrete types from the stdlib and ancestor
  packages.
  - e.g., `node/pipeline` can depend on `node`, but it shouldn't depend on
    `node/pipeline/handlers`
  - The exception to this is the root package, which is for general dependency
    injection wiring.
- Each package should only depend on abstract types from its descendant packages.
  - e.g., `node/pipeline/handlers` can depend on an abstract type that is
    implemented by a concrete type in `node/pipeline/handlers/store`
- `interfaces`'s should be used liberally to remove complexity from a package.
  - e.g., consider the use of `builder.Currenter` and `builder.Callback`

## C++

This project follows the
[Google C++ Style Guide](https://google.github.io/styleguide/cppguide.html),
with one exception: exceptions. See below for more details.

### Naming

#### Classes

When a class inherits from another class, and is the only class in source to
inherit from the other class, then it should be named with the name of the class
that it inherits from plus an `Impl` on the end.

E.g., say there is a base class called `Tuna`. If there is only one class in source
to inherit from the other class, then it should be named `TunaImpl`. If there is
a class in test that inherits from this base class as a mock, it should be named
`MockTuna`.

Otherwise, if there are multiple class that inherit from one base class, those
derived classes can be named whatever, preferably indicating that they derive from
the base class.