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

func TestPredecessorSuccessor(t *testing.T) {
	cases := []struct {
		input string
		index int // index in root that we want to find predecessor for
		pred  int // key of successor
		succ  int // key of scucessor
	}{
		{input: "(14(0)(23)(89))", index: 1, pred: 3, succ: 8},
		{input: "(4(2(1)(3))(68(5)(7)(9)))", index: 0, pred: 3, succ: 5},
		{input: "(P(CGM(AB)(DEF)(JKL)(NO))(TX(QRS)(UV)(YZ)))", index: 0, pred: 'O', succ: 'Q'},
		{input: "(TX(QRS)(UV)(YZ))", index: 1, pred: 'V', succ: 'Y'},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s/%d", tc.input, tc.index), func(t *testing.T) {
			tree := FromString(2, tc.input, os.Stderr)
			if leaf, i := tree.predecessor(tree.Root, tc.index); leaf.Keys[i] != tc.pred {
				got := leaf.Keys[i]
				t.Fatalf("predecessor mismatch; want=%q, got=%q", tc.pred, got)
			}
			if leaf, i := tree.successor(tree.Root, tc.index); leaf.Keys[i] != tc.succ {
				got := leaf.Keys[i]
				t.Fatalf("successor mismatch; want=%d, got=%d", tc.succ, got)
			}
		})
	}
}
func XXFoo() {
}

func TestMerge(t *testing.T) {
	cases := []struct {
		input string
		index int
		want  string
	}{
		{input: "(25(1)(4)(6))", index: 0, want: "(5(124)(6))"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			tree := FromString(2, tc.input, os.Stderr)
			tree.merge(tree.Root, tc.index)
			expectTree(t, tree, tc.want)
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		input string
		key   int
		want  string
	}{
		{input: "(123)", key: 2, want: "(13)"},
		{input: "(2(1)(34))", key: 4, want: "(2(1)(3))"},
		{input: "(3(12)(45))", key: 3, want: "(2(1)(45))"},   // case 2a
		{input: "(3(1)(45))", key: 3, want: "(4(1)(5))"},     // case 2b
		{input: "(25(1)(4)(6))", key: 2, want: "(5(14)(6))"}, // case 2c
		{input: "(5(4)(678))", key: 4, want: "(6(5)(78))"},   // case 3a
		{input: "(4(123)(5))", key: 5, want: "(3(12)(4))"},   // case 3b
		{input: "(24(1)(3)(5))", key: 3, want: "(2(1)(45))"}, // case 3b
		{input: "(2(1)(3))", key: 2, want: "(13)"},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s/%d", tc.input, tc.key), func(t *testing.T) {
			tree := FromString(2, tc.input, os.Stderr)
			// Graphviz(tree, "/tmp/xxyy.png")
			tree.Delete(tc.key)
			expectTree(t, tree, tc.want)
		})
	}

	// we apply these operations in sequence. Figure 18.8 in Cormen, p. 514
	t.Run("Cormen", func(t *testing.T) {
		cases := []struct {
			key  int
			want string
		}{
			{key: 'F', want: "(P(CGM(AB)(DE)(JKL)(NO))(TX(QRS)(UV)(YZ)))"},
			{key: 'M', want: "(P(CGL(AB)(DE)(JK)(NO))(TX(QRS)(UV)(YZ)))"},
			{key: 'G', want: "(P(CL(AB)(DEJK)(NO))(TX(QRS)(UV)(YZ)))"},
			{key: 'D', want: "(CLPTX(AB)(EJK)(NO)(QRS)(UV)(YZ))"},
			{key: 'B', want: "(ELPTX(AC)(JK)(NO)(QRS)(UV)(YZ))"},
		}
		tree := FromString(3, "(P(CGM(AB)(DEF)(JKL)(NO))(TX(QRS)(UV)(YZ)))", os.Stderr)
		for _, tc := range cases {
			tree.Delete(tc.key)
			expectTree(t, tree, tc.want)
		}
	})
}

func expectTree(t *testing.T, got *BTree, want string) {
	t.Helper()
	gotStr := got.String(got.Root)
	if gotStr != want {
		t.Fatalf("unexpected tree structure;\nwant= %s\ngot = %s", want, gotStr)
	}
}
