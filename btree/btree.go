package btree

import (
	"io"
	"log/slog"
)

type BTree struct {
	log  *slog.Logger
	Root *Node
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
	b.log.Info("Allocate-Node")
	return &Node{}
}
func (b *BTree) read(n *Node, i int) *Node {
	c := n.Children[i]
	b.log.Info("Disk read", "node", c)
	return c
}
func (b *BTree) write(n *Node) *Node {
	b.log.Info("Disk write", "node", n)
	return n
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

type Node struct {
	Leaf     bool
	Keys     []int
	Children []*Node
}

func (n *Node) String() string {
	var letters string
	for _, k := range n.Keys {
		letters += string(rune(k))
	}
	return letters
}
