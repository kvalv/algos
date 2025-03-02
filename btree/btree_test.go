package btree

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	// based on Figure 18.1 in Cormen, p. 498
	char := func(c rune) int { return int(c) }

	BC := Node{Keys: []int{char('B'), char('C')}, Leaf: true}
	FG := Node{Keys: []int{char('F'), char('G')}, Leaf: true}
	JKL := Node{Keys: []int{char('J'), char('K'), char('L')}, Leaf: true}

	NP := Node{Keys: []int{char('N'), char('P')}, Leaf: true}
	RS := Node{Keys: []int{char('R'), char('S')}, Leaf: true}
	VW := Node{Keys: []int{char('V'), char('W')}, Leaf: true}
	YZ := Node{Keys: []int{char('Y'), char('Z')}, Leaf: true}

	DH := Node{Keys: []int{char('D'), char('H')}, Children: []*Node{&BC, &FG, &JKL}}
	QTX := Node{Keys: []int{char('Q'), char('T'), char('X')}, Children: []*Node{&NP, &RS, &VW, &YZ}}

	M := Node{Keys: []int{char('M')}, Children: []*Node{&DH, &QTX}}

	btree := NewWithRoot(2, &M, os.Stderr)

	cases := []struct {
		key   rune
		leaf  *Node
		index int
	}{
		{key: 'R', leaf: &RS, index: 0},
		{key: 'T', leaf: &QTX, index: 1},
		{key: '1', leaf: nil, index: 0},
		{key: 'B', leaf: &BC, index: 0},
		{key: 'Z', leaf: &YZ, index: 1},
		{key: 'M', leaf: &M, index: 0},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%c", tc.key), func(t *testing.T) {
			leaf, index := btree.Search(btree.Root, char(tc.key))

			if tc.leaf != leaf {
				t.Errorf("Node mismatch; want=%v, got=%v", tc.leaf, leaf)
			}
			if tc.index != index {
				t.Fatalf("index mismatch; want=%v, got=%v", tc.index, index)
			}
		})
	}
}

func TestKeys(t *testing.T) {
	// based on Figure 18.1 in Cormen, p. 498
	char := func(c rune) int { return int(c) }

	BC := Node{Keys: []int{char('B'), char('C')}, Leaf: true}
	FG := Node{Keys: []int{char('F'), char('G')}, Leaf: true}
	JKL := Node{Keys: []int{char('J'), char('K'), char('L')}, Leaf: true}

	NP := Node{Keys: []int{char('N'), char('P')}, Leaf: true}
	RS := Node{Keys: []int{char('R'), char('S')}, Leaf: true}
	VW := Node{Keys: []int{char('V'), char('W')}, Leaf: true}
	YZ := Node{Keys: []int{char('Y'), char('Z')}, Leaf: true}

	DH := Node{Keys: []int{char('D'), char('H')}, Children: []*Node{&BC, &FG, &JKL}}
	QTX := Node{Keys: []int{char('Q'), char('T'), char('X')}, Children: []*Node{&NP, &RS, &VW, &YZ}}

	M := Node{Keys: []int{char('M')}, Children: []*Node{&DH, &QTX}}

	btree := NewWithRoot(2, &M, os.Stderr)
	got := btree.Keys()

	var want []int
	for _, c := range "BCDFGHJKLMNPQRSTVWXYZ" {
		want = append(want, int(c))
	}

	if len(got) != len(want) {
		t.Fatalf("length mismatch; want=%d, got=%d", len(want), len(got))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("mismatch at index %d; want=%d, got=%d", i, want[i], got[i])
		}
	}
}

func TestSplit(t *testing.T) {
	right := &Node{
		Keys: []int{5, 6, 7},
		Leaf: true,
	}
	left := &Node{
		Keys: []int{1, 2},
		Leaf: true,
	}
	root := &Node{
		Keys:     []int{3},
		Children: []*Node{left, right},
	}
	tree := NewWithRoot(2, root, os.Stderr)

	expectTree(t, tree, "(3(12)(567))")
	tree.SplitChild(root, 1)
	expectTree(t, tree, "(36(12)(5)(7))")
}

func TestSplitRoot(t *testing.T) {
	cases := []struct {
		keys []int
		want string
	}{
		{keys: []int{1, 2, 3}, want: "(2(1)(3))"},
		{keys: []int{1, 2, 3, 4}, want: "(3(12)(4))"},
		{keys: []int{1, 2, 3, 4, 5}, want: "(3(12)(45))"},
		{keys: []int{2, 2, 2}, want: "(2(2)(2))"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.keys), func(t *testing.T) {
			root := &Node{Keys: tc.keys, Leaf: true}
			tree := NewWithRoot(3, root, os.Stderr)
			tree.dbg = true
			tree.splitRoot()
			expectTree(t, tree, tc.want)

			tree2 := FromString(2, tc.want, io.Discard)
			expectTree(t, tree2, tc.want)
		})
	}
}

func TestSplitRootV2(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			input: "(246(1)(3)(5)(78))",
			want:  "(4(2(1)(3))(6(5)(78)))",
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			tree := FromString(2, tc.input, os.Stderr)
			tree.splitRoot()
			Graphviz(tree, "/tmp/btreexx.png")
			expectTree(t, tree, tc.want)
		})
	}
}

func TestInsert(t *testing.T) {
	cases := []struct {
		n    int
		keys []int
		want string
	}{
		{n: 2, keys: []int{1, 2, 3, 4}, want: "(2(1)(34))"},
		{n: 2, keys: []int{1, 2, 3, 4, 5, 6, 7}, want: "(24(1)(3)(567))"},
		{n: 2, keys: []int{1, 2, 3, 4, 5, 6, 7, 8}, want: "(246(1)(3)(5)(78))"},
		{n: 2, keys: []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, want: "(4(2(1)(3))(6(5)(789)))"},
		{n: 2, keys: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, want: "(4(2(1)(3))(68(5)(7)(910)))"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%d", tc.n), func(t *testing.T) {
			tree := New(tc.n, os.Stderr)
			for _, key := range tc.keys {
				tree.Insert(key)
			}
			// Graphviz(tree, "/tmp/xx.png")
			expectTree(t, tree, tc.want)

			tree2 := FromString(2, tc.want, io.Discard)
			expectTree(t, tree2, tc.want)

		})
	}

}

func expectTree(t *testing.T, got *BTree, want string) {
	t.Helper()
	gotStr := got.String(got.Root)
	if gotStr != want {
		t.Fatalf("unexpected tree structure;\nwant= %s\ngot = %s", want, gotStr)
	}
}
