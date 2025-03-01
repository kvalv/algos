package btree

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type BTree struct {
	log   *slog.Logger
	Root  *Node
	stats Stats
}

// Keep statistics about read / write access
type Stats struct {
	Reads, Writes int
}

func New(root *Node, w io.Writer) *BTree {
	return &BTree{
		log:  NewLogger(w),
		Root: root,
	}
}

func Create(w io.Writer) *BTree {
	b := &BTree{log: NewLogger(w)}
	x := b.allocate()
	x.Leaf = true
	b.write(x)
	b.Root = x
	return b
}

func (b *BTree) allocate() *Node {
	b.log.Debug("Allocate-Node")
	return &Node{}
}
func (b *BTree) read(n *Node, i int) *Node {
	c := n.Children[i]
	_, med := n.median()
	b.log.Debug("Disk read", "node", keyString(med))
	b.stats.Reads++
	return c
}
func (b *BTree) write(n *Node) *Node {
	_, med := n.median()
	b.log.Debug("Disk write", "node", keyString(med))
	b.stats.Writes++
	return n
}

// x.Children[i] is assumed full; x is assumed non-full. We split the child and
// put the median key into x
func (b *BTree) SplitChild(x *Node, i int) {
	// c := b.read(x, i)
	// idx, k := c.median()
}

func (b *BTree) Search(n *Node, key int) (*Node, int) {
	for i, k := range n.Keys {
		if k == key {
			return n, i
		} else if k > key {
			if n.Leaf {
				return nil, 0
			}
			// return b.sea
			// c := n.Children[i] // disk read
			c := b.read(n, i)
			return b.Search(c, key)
		}
	}
	if n.Leaf {
		return nil, 0
	}
	c := n.Children[len(n.Keys)] // disk read here
	return b.Search(c, key)
}

func (b *BTree) Walk(n *Node, f func(key int)) {
	if n == nil {
		return
	}
	for i, key := range n.Keys {
		if !n.Leaf {
			c := b.read(n, i)
			b.Walk(c, f)
		}
		f(key)
	}
	if !n.Leaf {
		c := b.read(n, len(n.Keys))
		b.Walk(c, f)
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

func (b *BTree) String(n *Node) string {
	var s strings.Builder
	s.WriteString("(")
	for _, k := range n.Keys {
		s.WriteString(keyString(k))
	}
	for i := range n.Children {
		c := b.read(n, i)
		fmt.Fprintf(&s, "%s", b.String(c))
	}
	s.WriteString(")")

	return s.String()
}

func (n *Node) median() (index int, key int) {
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
