package bplus

import (
	"fmt"
	"strings"
)

func (n *Node) String() string {
	if n == nil {
		return "nil"
	}
	var s strings.Builder
	s.WriteString("(")
	for _, k := range n.Keys {
		s.WriteString(keyString(k))
	}
	s.WriteString(")")
	return s.String()
}

func (T *BTree) String(n *Node) string {
	if n == nil {
		panic("BTree.String(): n is nil")
	}
	var s strings.Builder
	s.WriteString("(")
	for _, k := range n.Keys {
		s.WriteString(keyString(k))
	}
	if !n.Leaf { // I think we're using .Children for pointers, too...
		for i := range n.Children {
			c := T.read(n, i)
			fmt.Fprintf(&s, "%s", T.String(c))
		}
	}
	s.WriteString(")")

	return s.String()
}
