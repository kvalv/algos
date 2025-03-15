package bplus

import (
	"fmt"
	"log/slog"
	"slices"
)

type BTree struct {
	// t = n - 1
	n         int // n pointers, n-1 keys
	log       *slog.Logger
	dbg       bool
	Root      *Node
	pageCache *PageCache
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
	id := n.Children[i]
	return T.pageCache.Read(id)
}
func (T *BTree) write(n *Node) *Node {
	return T.pageCache.Write(n)
}
func (T *BTree) allocate() *Node {
	return T.pageCache.Allocate()
}

func (T *BTree) WalkNodes(n *Node, f func(n *Node)) {
	if n == nil {
		return
	}
	f(n)
	for _, id := range n.Children {
		c := T.pageCache.Read(id)
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

// returns nil if node does not have any keys, or the key is greater
// than all keys in this set, in which case - consider calling T.lastChild
func (T *BTree) insertionIndex(node *Node, key int) *int {
	if len(node.Keys) == 0 || key > node.Keys[len(node.Keys)-1] {
		return nil
	}
	// otherwise we knwo for sure there's at least one key that is greater
	for i, k := range node.Keys {
		if k >= key {
			return &i
		}
	}
	panic("unreachable")
}

func (T *BTree) lastChild(N *Node) *Node {
	length := len(N.Children)
	if length == 0 {
		return nil
	}
	pageID := N.Children[length-1]
	return T.pageCache.Read(pageID)
}

func (T *BTree) Find(key int) *Match {
	C := T.Root
	var count int
	for !C.Leaf {
		count++
		if count > 100 {
			panic("infinite loop")
		}
		i := T.insertionIndex(C, key)
		if i == nil {
			C = T.lastChild(C)
			continue
		}
		if C.Keys[*i] == key {
			C = T.read(C, *i+1)
		} else {
			C = T.read(C, *i)
		}
	}
	i := T.insertionIndex(C, key)
	if i == nil {
		return nil
	}
	return &Match{C, *i}
}

func (T *BTree) Range(key, upper int) RangeIterator {
	C := T.Root
	for !C.Leaf {
		i := T.insertionIndex(C, key)
		if i == nil {
			C = T.lastChild(C)
			continue
		}
		if C.Keys[*i] == key {
			C = T.read(C, *i+1)
		} else {
			C = T.read(C, *i)
		}
	}
	i := T.insertionIndex(C, key)
	if i == nil {
		return NewEmptyIterator[Match]()
	}

	j := *i
	return NewIterator(func() *Match {
		m := Match{Index: j, Node: C}
		if j >= len(C.Keys) {
			if C.RightSibling == nil {
				return nil
			}
			// otherwise visit next sibling
			C = T.pageCache.Read(*C.RightSibling)
			j = 0
			m = Match{Index: j, Node: C}
		}
		if C.Keys[j] >= upper {
			return nil
		}
		j++
		return &m
	})
}

func (T *BTree) Insert(key int, value PageID) {
	node := T.Root

	var parent *Node
	for !node.Leaf {
		parent = node
		if i := T.insertionIndex(node, key); i != nil {
			node = T.read(node, *i)
		} else {
			node = T.lastChild(node)
		}
	}

	// now we are at the leaf, and we can insert...
	i := T.insertionIndex(node, key)
	if i == nil {
		node.Keys = append(node.Keys, key)
		node.Children = append(node.Pointers, value)
	} else {
		node.Keys = slices.Insert(node.Keys, *i, key)
		node.Children = slices.Insert(node.Pointers, *i, value)
	}

	// split step ... if
	if len(node.Keys) <= T.n {
		return // we're done!
	}

	// otherwise we need to split
	j := (T.n + 1) / 2 // ceil[n/2]
	right := T.SplitLeaf(node, j)

	if parent == nil {
		// create a new root
		root := T.allocate()
		root.Keys = []int{right.MinKey()}
		root.Children = []PageID{node.PageID, right.PageID}
		T.Root = root
	} else {
		panic("TODO: add key to existing parent")
	}

}

// Splits current node at index i, returning the new node, residing on the right side
func (T *BTree) SplitLeaf(node *Node, i int) *Node {
	right := T.pageCache.Allocate()

	right.Pointers = node.Pointers[i:]
	right.Keys = node.Keys[i:]
	right.RightSibling = node.RightSibling

	node.Pointers = node.Pointers[:i]
	node.Keys = node.Keys[:i]
	tmp := right.PageID
	node.RightSibling = &tmp

	return right
}
