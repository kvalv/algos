package btree

import (
	"fmt"
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

	btree := New(&M, os.Stderr)

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

	btree := New(&M, os.Stderr)
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
	// t=2 is the simplest type, then t-1 to 2t-1 keys -> 1 to 3 keys -> 2 to 4 children. It is full if 4 children

	right := &Node{
		Keys: []int{5, 6, 7},
	}
	left := &Node{
		Keys: []int{1, 2},
	}
	root := &Node{
		Keys:     []int{3},
		Children: []*Node{left, right},
	}
	tree := New(root, os.Stderr)

	if want, got := "(3(12)(567))", tree.String(root); want != got {
		t.Fatalf("unexpected tree structure;\nwant= %s\ngot = %s", want, got)
	}

	tree.SplitChild(root, 1) //
	want := "(36(12)(5)(7))"
	if got := tree.String(root); want != got {
		t.Fatalf("unexpected tree structure;\nwant= %s\ngot = %s", want, got)
	}

	// want := "(36(12)(5)(7))"

}
