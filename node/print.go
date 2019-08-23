package node

import (
	"bytes"
	"fmt"
	"io"
)

// String returns a text-based representation of the Node graph.
func String(n *Node) string {
	buf := bytes.NewBuffer([]byte{})
	streeng(n, buf, "", 10)
	return buf.String()
}

func streeng(n *Node, buf io.Writer, prefix string, depth int) {
	if depth == 0 {
		return
	}

	fmt.Fprintf(buf, "%s%s\n", prefix, n.Name)
	for _, dN := range n.Dependencies {
		streeng(dN, buf, prefix+". ", depth-1)
	}
}
