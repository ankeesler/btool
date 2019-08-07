# Coding Standard

- `node.Node` variables should always end in an "n" or an "N"
  - e.g., `dN` (for a `node.Node.Dependency`), `n` (for a `node.Node`)
  - e.g., `oN` (for an object file node), `oC` (for a c file node)
- `[]*node.Node` variables should always end in "nodes" or "Nodes"
  - e.g., `depNodes` (for the `node.Node`'s that describe a dep)