# Coding Standard

This project follows the
[Google C++ Style Guide](https://google.github.io/styleguide/cppguide.html),
with one exception: exceptions. See below for more details.

## Naming

### Classes

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