package rb

import (
	"testing"
)

func TestRotate(t *testing.T) {
	var orig, rotated Tree
	orig.Insert(5, 2, 4, 1, 3)

	root := &Node{Key: 2, Color: BLACK}
	left := &Node{Key: 1, Parent: root}
	root.Left = left
	right := &Node{Key: 4, Parent: root}
	root.Right = right
	right.Left = &Node{Key: 3, Parent: right}
	right.Right = &Node{Key: 5, Parent: right}

	rotated = Tree{Root: root}

	wantOrder := []int{1, 2, 3, 4, 5}
	expectOrder(t, &orig, wantOrder)

	x := rotated.Find(2)

	rotated.LeftRotate(x)
	rotated.RightRotate(rotated.Find(4))
	expectOrder(t, &rotated, wantOrder)
	// Graphviz(rotated, "/tmp/origxy.png")
}

func TestInsert(t *testing.T) {
	var tree Tree
	tree.Insert(1, 2, 4, 5, 8, 7, 11, 14, 15)
	Graphviz(tree, "/tmp/origxc.png")
}

func expectOrder(t *testing.T, tree *Tree, want []int) {
	t.Helper()
	var got []int
	tree.Walk(func(n *Node) {
		got = append(got, n.Key)
	})
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}

}
