# Coding Standard

## Directory Structure

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