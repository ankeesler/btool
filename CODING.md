# Coding Standard

## Coding Format

- `node.Node` variables should always end in an "n" or an "N"
  - e.g., `dN` (for a `node.Node.Dependency`), `n` (for a `node.Node`)
  - e.g., `oN` (for an object file node), `oC` (for a c file node)
- `[]*node.Node` variables should always end in "nodes" or "Nodes"
  - e.g., `depNodes` (for the `node.Node`'s that describe a dep)

## Directory Structure

- Each abstraction should be a subdirectory the abstraction that it is built
  on top of
  - e.g., `pipeline` is a subdirectory of `node`
- Any abstraction that is based on 2 or more other abstractions should be a
  subdirectory of the highest possible abstraction
  - e.g., `handlers.NewObject()` depepds on `resolvers.NewCompile()` and `node.Node`,
    therefore it should be a subdirectory of `node`
  - TODO: provide a `ResolverFactory` to the pipeline so we don't have the above
    abstraction issue