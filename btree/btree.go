package btree

import (
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strings"
)

type BTree struct {
	n     int
	log   *slog.Logger
	Root  *Node
	stats Stats
	dbg   bool
}

// Keep statistics about read / write access
type Stats struct {
	Reads, Writes int
}

func NewWithRoot(n int, root *Node, w io.Writer) *BTree {
	return &BTree{
		n:    n,
		log:  NewLogger(w),
		Root: root,
	}
}

// Helper function, primarily for writing test cases
func FromString(n int, input string, w io.Writer) *BTree {
	T := New(n, w)

	var stack []*Node
	top := func() *Node {
		if len(stack) == 0 {
			return nil
		}
		return stack[len(stack)-1]
	}
	var root *Node
	pop := func() {
		if len(stack) == 0 {
			return
		}
		if len(stack) == 1 {
			root = stack[0]
		}
		stack = stack[:len(stack)-1]
	}

	toKey := func(c rune) int {
		// if 0123...f -> parse as hex {
		if c >= '0' && c <= '9' {
			return int(c - '0')
		}
		return int(c)
	}

	for _, c := range input {
		switch c {
		case '(':
			tmp := &Node{}
			p := top()
			if p != nil {
				p.Children = append(p.Children, tmp)
			}
			stack = append(stack, tmp)
		case ')':
			pop()
		default:
			tmp := top()
			digit := toKey(c)
			tmp.Keys = append(tmp.Keys, digit)
		}
	}

	T.Root = root
	return T
}

func New(n int, w io.Writer) *BTree {
	b := &BTree{
		n:   n,
		log: NewLogger(w),
	}
	x := b.allocate()
	x.Leaf = true
	b.write(x)
	b.Root = x
	return b
}

func (T *BTree) allocate() *Node {
	T.log.Debug("Allocate-Node")
	return &Node{}
}
func (T *BTree) read(n *Node, i int) *Node {
	c := n.Children[i]
	_, med := n.median()
	T.log.Debug("Disk read", "node", keyString(med))
	T.stats.Reads++
	return c
}
func (T *BTree) write(n *Node) *Node {
	_, med := n.median()
	T.log.Debug("Disk write", "node", keyString(med))
	T.stats.Writes++
	return n
}

func (T *BTree) validateTree() {
	T.WalkNodes(T.Root, func(n *Node) {
		if len(n.Keys) > 2*T.n-1 {
			panic(fmt.Sprintf("node %s has %d keys", n, len(n.Keys)))
		}
		if !n.Leaf && len(n.Keys)+1 != len(n.Children) {
			panic(fmt.Sprintf("node %s has %d keys but %d children", n, len(n.Keys), len(n.Children)))
		}
	})
}

// x.Children[i] is assumed full; x is assumed non-full. We split the child and
// put the median key into x
func (T *BTree) SplitChild(x *Node, i int) int {
	y := T.read(x, i)
	z := T.allocate()
	z.Leaf = y.Leaf

	medianIndex := len(y.Keys) / 2
	key := y.Keys[medianIndex]
	T.log.Info(
		"splitChild",
		"y", y,
		"medianKey", key,
		"children", y.Children,
		"leaf", y.Leaf,
	)
	z.Keys = y.Keys[medianIndex+1:]
	if !y.Leaf && len(y.Children) > 0 {
		z.Children = y.Children[medianIndex+1:]
		y.Children = y.Children[:medianIndex+1]
	}

	y.Keys = y.Keys[:medianIndex]

	T.log.Info("x.keys.i", "keys", x.Keys, "i", i)
	x.Keys = slices.Insert(x.Keys, i, key)
	x.Children = slices.Insert(x.Children, i+1, z)

	T.write(x)
	T.write(y)
	T.write(z)
	if T.dbg {
		if !y.Leaf {
			if len(y.Keys)+1 != len(y.Children) {
				panic(fmt.Sprintf("y has %d keys but %d children keys=%v children=%v", len(y.Keys), len(y.Children), y.Keys, y.Children))
			}
			if len(z.Keys)+1 != len(z.Children) {
				panic(fmt.Sprintf("z has %d keys but %d children", len(z.Keys), len(z.Children)))
			}
		}
		if len(x.Keys)+1 != len(x.Children) {
			panic(fmt.Sprintf("x has %d keys but %d children", len(x.Keys), len(x.Children)))
		}
		T.validateTree()
	}

	return key
}

func (T *BTree) full(x *Node) bool     { return len(x.Keys) == 2*T.n-1 }
func (T *BTree) starving(x *Node) bool { return len(x.Keys) == T.n-1 } // bad name, TODO

func (T *BTree) Insert(key int) {
	x := T.Root
	if T.full(T.Root) {
		x = T.splitRoot()
	}
	T.insertNonFull(x, key)

}

// assume non-full when this method is called.
func (T *BTree) insertNonFull(x *Node, key int) {
	i := x.indexFor(key)

	if x.Leaf {
		x.Keys = slices.Insert(x.Keys, i, key)
		return
	}
	// else it's not a leaf, so we check if it's full or not
	c := T.read(x, i)
	if T.full(c) {
		med := T.SplitChild(x, i)
		if key > med {
			c = T.read(x, i+1)
		}
	}
	T.insertNonFull(c, key)
}

// returns the new root
func (T *BTree) splitRoot() *Node {
	s := T.allocate()
	// s.Leaf = false
	s.Children = []*Node{T.Root}
	T.Root = s
	T.SplitChild(s, 0)
	s.Leaf = false
	return s
}

func (T *BTree) Search(n *Node, key int) (*Node, int) {
	for i, k := range n.Keys {
		if k == key {
			return n, i
		} else if k > key {
			if n.Leaf {
				return nil, 0
			}
			// return b.sea
			// c := n.Children[i] // disk read
			c := T.read(n, i)
			return T.Search(c, key)
		}
	}
	if n.Leaf {
		return nil, 0
	}
	c := n.Children[len(n.Keys)] // disk read here
	return T.Search(c, key)
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
	Leaf     bool
	Keys     []int
	Children []*Node
}

// Index of first element that is smaller than key.
// If key is greater than all elements, len(Node.Keys) is returned.
// So if n=2, then 3 is returned.
func (n *Node) indexFor(key int) int {
	for i, k := range n.Keys {
		if key < k {
			return i
		}
	}
	return len(n.Keys)
}

func (n *Node) String() string {
	var s strings.Builder
	s.WriteString("(")
	for _, k := range n.Keys {
		s.WriteString(keyString(k))
	}
	s.WriteString(")")
	return s.String()
}

func (T *BTree) String(n *Node) string {
	var s strings.Builder
	s.WriteString("(")
	for _, k := range n.Keys {
		s.WriteString(keyString(k))
	}
	for i := range n.Children {
		c := T.read(n, i)
		fmt.Fprintf(&s, "%s", T.String(c))
	}
	s.WriteString(")")

	return s.String()
}

func (n *Node) median() (index int, key int) {
	if len(n.Keys) == 0 {
		return 0, 0
	}

	index = len(n.Keys) / 2
	key = n.Keys[index]
	return
}

func keyString(k int) string {
	if k >= 'A' && k <= 'Z' || k >= 'a' && k <= 'z' {
		return fmt.Sprintf("%c", k)
	}
	return fmt.Sprintf("%d", k)
}
