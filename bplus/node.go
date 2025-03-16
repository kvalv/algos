package bplus

import "strings"

type Node struct {
	PageID
	Keys []int
	Leaf bool

	RightSibling *PageID

	// leaf: has N-1 keys and N pointers
	// For leaf, the last pointer points to sibling node (next) - not back
	Children []PageID // child nodes
	Values   []PageID // things we point to
}

func (n *Node) median() (index int, key int) {
	if len(n.Keys) == 0 {
		return 0, 0
	}

	index = len(n.Keys) / 2
	key = n.Keys[index]
	return
}

func (n *Node) MinKey() int {
	if len(n.Keys) == 0 {
		panic("FirstKey: Node has no keys")
	}
	return n.Keys[0]
}

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
