package rb

import "fmt"

/*
1. every node marked red or black
2. leaf (nil) nodes are black
3. red nodes do not have red childnen
4. Every path from a particular node to a NIL node must go through the same number of black nodes

Properties
Allows efficient traversal,

*/

// red has black children
// root is black; NIL's black
// all paths from node to a leaf contain same number of black nodes

type Node struct {
	Key                 int
	Left, Right, Parent *Node
	Color
}
type Color int

func (n *Node) String() string { return fmt.Sprintf("%d", n.Key) }

type Tree struct {
	Root *Node
}

const (
	RED Color = iota
	BLACK
)

func (t *Tree) LeftRotate(x *Node) {
	y := x.Right

	x.Right = y.Left
	if y.Left != nil {
		y.Left.Parent = x
	}

	y.Parent = x.Parent

	y.Left = x
	if t.Root == x {
		t.Root = y
	} else if x.Parent.Left == x {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}
	x.Parent = y
}

func (t *Tree) RightRotate(x *Node) {
	y := x.Left

	x.Left = y.Right
	if y.Right != nil {
		y.Right.Parent = x
	}

	y.Parent = x.Parent

	y.Right = x
	if t.Root == x {
		t.Root = y
	} else if x.Parent.Left == x {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}
	x.Parent = y
}

func (t *Tree) Insert(key ...int) *Tree {
	for _, k := range key {
		t.insert(k)
	}
	return t
}
func (t *Tree) Find(key int) *Node {
	n := t.Root
	for n != nil && n.Key != key {
		if n.Key > key {
			n = n.Left
		} else {
			n = n.Right
		}
	}
	return n
}

func (t *Tree) Walk(f func(n *Node)) {
	t.Root.walk(f)
}
func (n *Node) walk(f func(n *Node)) {
	if n == nil {
		return
	}
	fmt.Printf("walk: %d\n", n.Key)
	n.Left.walk(f)
	f(n)
	n.Right.walk(f)
}

func (t *Tree) insert(key int) {
	// start naive for now
	n := t.Root
	var par *Node

	for n != nil {
		par = n
		if n.Key > key {
			n = n.Left
		} else {
			n = n.Right
		}
	}
	x := &Node{Key: key, Parent: par}
	if par == nil {
		t.Root = x
	} else if par.Key > key {
		par.Left = x
	} else {
		par.Right = x
	}
}
