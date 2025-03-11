package bplus

import (
	"fmt"
	"log/slog"
)

type BTree struct {
	// t = n - 1
	n    int // n pointers, n-1 keys
	log  *slog.Logger
	dbg  bool
	Root *Node
}

func (T *BTree) validate() {
	if !T.dbg {
		return
	}
	T.WalkNodes(T.Root, func(n *Node) {
		if len(n.Keys) > 2*T.n-1 {
			panic(fmt.Sprintf("node %s has %d keys", n, len(n.Keys)))
		}
		if !n.Leaf && len(n.Keys)+1 != len(n.Children) {
			panic(fmt.Sprintf("BTree violation: node %s has %d and %d children - expected %d children",
				n,
				len(n.Keys),
				len(n.Children),
				len(n.Keys)+1),
			)
		}
		if len(n.Children) > 0 && n.Leaf {
			panic("Node is a leaf - but has children")
		}
		if len(n.Children) == 0 && !n.Leaf {
			panic("Node is not a leaf, but has children")
		}
	})
}
func (T *BTree) read(n *Node, i int) *Node {
	if len(n.Children) < i-1 || i < 0 {
		return nil
	}
	c := n.Children[i]
	T.log.Debug("Disk read", "node", c.String())
	return c
}
func (T *BTree) write(n *Node) *Node {
	_, med := n.median()
	T.log.Debug("Disk write", "node", keyString(med))
	return n
}
func (n *Node) median() (index int, key int) {
	if len(n.Keys) == 0 {
		return 0, 0
	}

	index = len(n.Keys) / 2
	key = n.Keys[index]
	return
}
func (T *BTree) allocate() *Node {
	T.log.Debug("Allocate-Node")
	return &Node{}
}

func (T *BTree) WalkNodes(n *Node, f func(n *Node)) {
	if n == nil {
		return
	}
	f(n)
	for _, c := range n.Children {
		T.WalkNodes(c, f)
	}
}
func (T *BTree) Walk(n *Node, f func(key int)) {
	if n == nil {
		return
	}
	for i, key := range n.Keys {
		if !n.Leaf {
			c := T.read(n, i)
			T.Walk(c, f)
		}
		f(key)
	}
	if !n.Leaf {
		c := T.read(n, len(n.Keys))
		T.Walk(c, f)
	}
}
func (b *BTree) Keys() []int {
	var res []int
	b.Walk(b.Root, func(key int) {
		res = append(res, key)
	})
	return res
}

type Node struct {
	Keys []int
	Leaf bool

	// leaf: has N-1 keys and N pointers
	// For leaf, the last pointer points to sibling node (next) - not back
	Children []*Node // ??
	Pointers []PageID
}

func (T *BTree) Find(key int) (*Node, int) {
	C := T.Root
	for !C.Leaf {
		return C.Children[1], 1
	}
	return nil, 0
}
