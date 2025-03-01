// based on https://web.stanford.edu/class/archive/cs/cs161/cs161.1168/lecture8.pdf

package btree

import (
	"fmt"
	"io"
)

type Node struct {
	key                 int
	Parent, Left, Right *Node
}

type BTree struct {
	Root *Node
}

func New() *BTree {
	return &BTree{}
}

func (b *BTree) Insert(key ...int) *BTree {
	for _, k := range key {
		b.insert(k)
	}
	return b
}

func (b *BTree) insert(key int) {
	var par *Node
	curr := b.Root
	for curr != nil {
		par = curr
		if curr.key > key {
			curr = curr.Left
		} else {
			curr = curr.Right
		}
	}
	if par == nil {
		b.Root = &Node{key: key}
		return
	}
	if par.key > key {
		par.Left = &Node{key: key, Parent: par}
	} else {
		par.Right = &Node{key: key, Parent: par}
	}
}

func (b *BTree) Remove(key int) *BTree {
	z := b.Root.Find(key)
	if z == nil {
		return b // our job is done
	}

	par := z.Parent
	switch z.countDirectChildren() {
	case 0:
		// no children...
		if par.Left == z {
			par.Left = nil
		} else {
			par.Right = nil
		}
	case 1:
		b.transplant(z, z.first())
	case 2:
		s := z.Successor()
		if z.Right == s {
			// right successor should just be moved up; we do that by
			// moving the value over and then removing it
			z.Right = s.Right
			z.key = s.key

		} else {
			// case 2: successor is in right subtree, but not as direct
			// child. Get rid of successor and move key here. We know it's a leaf
			// so no need to worry about children
			b.Remove(s.key)
			z.key = s.key
		}
	default:
		panic("unexpected number of children")
	}
	return b
}

// replace subtree rooted at old with subtree rooted at new
func (b *BTree) transplant(old, new *Node) {
	// a node takes over for old. Hijacks its children. Assume it's child-less
	par := old.Parent
	if par == nil {
		b.Root = new
		return
	}
	if par.Left == old {
		par.Left = new
	} else if par.Right == old {
		par.Right = new
	}
	new.Parent = par
}

func (n *Node) Find(key int) *Node {
	for n != nil && n.key != key {
		if n.key > key {
			n = n.Left
		} else {
			n = n.Right
		}
	}
	return n
}

func (n *Node) Min() *Node {
	for n != nil && n.Left != nil {
		n = n.Left
	}
	return n
}

func (n *Node) Max() *Node {
	for n != nil && n.Right != nil {
		n = n.Right
	}
	return n
}

func (x *Node) Successor() *Node {
	if x.Right != nil {
		return x.Right.Min()
	}
	y := x.Parent
	for y != nil && x == y.Right {
		x = y
		y = y.Parent
	}
	return y
}

func (n *Node) Walk() []int {
	if n == nil {
		return nil
	}
	var res []int
	res = append(res, n.Left.Walk()...)
	res = append(res, n.key)
	res = append(res, n.Right.Walk()...)
	return res
}

func (n *Node) countDirectChildren() int {
	var count int
	if n.Left != nil {
		count++
	}
	if n.Right != nil {
		count++
	}
	return count
}

// return the first child of Node. panics if not found
func (n *Node) first() *Node {
	if n.Left != nil {
		return n.Left
	}
	if n.Right != nil {
		return n.Right
	}
	panic("no child")

}

func (n *Node) String() string {
	if n == nil {
		return ""
	}
	return fmt.Sprintf("(%d%s%s)", n.key, n.Left.String(), n.Right.String())
}
func (b *BTree) String() string {
	return b.Root.String()
}

func (b *BTree) Graphviz(w io.Writer) {
	fmt.Fprintf(w, "digraph G {\n %s\n}", b.Root.graphviz())
}
func (n *Node) graphviz() string {
	if n == nil {
		return ""
	}
	var res string
	if n.Left != nil {
		res += fmt.Sprintf("%d -> %d;\n", n.key, n.Left.key)
		res += n.Left.graphviz()
	}
	if n.Right != nil {
		res += fmt.Sprintf("%d -> %d;\n", n.key, n.Right.key)
		res += n.Right.graphviz()
	}
	return res
}
