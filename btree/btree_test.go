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
	// btree := &BTree{Root: &M}

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
