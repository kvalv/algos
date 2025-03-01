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

var (
	SENTINEL = &Node{Color: BLACK}
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
	if n == nil || n == SENTINEL {
		return
	}
	fmt.Printf("walk: %d\n", n.Key)
	n.Left.walk(f)
	f(n)
	n.Right.walk(f)
}

func (t *Tree) insert(key int) {
	// var par *Node
	par := SENTINEL
	x := t.Root
	z := &Node{Key: key, Color: RED, Left: SENTINEL, Right: SENTINEL, Parent: SENTINEL}

	for x != nil && x != SENTINEL {
		par = x
		if z.Key < x.Key {
			x = x.Left
		} else {
			x = x.Right
		}
	}
	z.Parent = par
	if par == SENTINEL {
		t.Root = z
	} else if par.Key > z.Key {
		par.Left = z
	} else {
		par.Right = z
	}
	// ^ ... so just regular BST insert, but with a correction step at the end.
	// no sentinel T.nil is used, but maybe we need to.
	t.InsertFixup(z)
}

func (t *Tree) InsertFixup(z *Node) {
	for z.Parent.Color == RED {
		if z.Parent == z.Parent.Parent.Left {
			y := z.Parent.Parent.Right // uncle
			if y.Color == RED {        // case 1: uncle is red; recolor father and uncle
				z.Parent.Color = BLACK
				y.Color = BLACK
				z.Parent.Parent.Color = RED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Right { // case 2
					z = z.Parent
					t.LeftRotate(z)
				}
				z.Parent.Color = BLACK
				z.Parent.Parent.Color = RED
				t.RightRotate(z.Parent.Parent)
			}
		} else {
			// mirrored version??
			y := z.Parent.Parent.Left // uncle
			if y.Color == RED {       // case 1: uncle is red; recolor father and uncle
				z.Parent.Color = BLACK
				y.Color = BLACK
				z.Parent.Parent.Color = RED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Left { // case 2
					z = z.Parent
					t.RightRotate(z)
				}
				z.Parent.Color = BLACK
				z.Parent.Parent.Color = RED
				t.LeftRotate(z.Parent.Parent)
			}
		}
	}
	t.Root.Color = BLACK
}

func (t *Tree) Transplant(u, v *Node) {
	if u.Parent == SENTINEL {
		t.Root = v
	} else if u.Parent.Left == u {
		u.Parent.Left = v
	} else {
		u.Parent.Right = v
	}
	v.Parent = u.Parent
}
