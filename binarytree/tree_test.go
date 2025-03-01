package btree

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestInsert(t *testing.T) {
	cases := []struct {
		desc string
		tree *BTree
		want string
	}{
		{
			desc: "simple",
			tree: New().Insert(2, 1, 3),
			want: "(2(1)(3))",
		},
		{
			desc: "balanced",
			tree: New().Insert(4, 2, 6, 1, 3, 5, 7),
			want: "(4(2(1)(3))(6(5)(7)))",
		},
		{
			desc: "left",
			tree: New().Insert(5, 4, 3, 2, 1),
			want: "(5(4(3(2(1)))))",
		},
		{
			desc: "right",
			tree: New().Insert(1, 2, 3, 4, 5),
			want: "(1(2(3(4(5)))))",
		},
		{
			desc: "unbalanced",
			tree: New().Insert(10, 5, 15, 3, 7, 18),
			want: "(10(5(3)(7))(15(18)))",
		},
		{
			desc: "root",
			tree: New().Insert(1),
			want: "(1)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.tree.String()
			if got != tc.want {
				t.Errorf("tree mismatch \nwant=%q, \ngot =%q", tc.want, got)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		desc string
		tree *BTree
		key  int
		want string
	}{
		{
			desc: "simple",
			tree: New().Insert(2, 1, 3),
			key:  1,
			want: "(2(3))",
		},
		{
			desc: "balanced",
			tree: New().Insert(4, 2, 6, 1, 3, 5, 7),
			key:  2,
			want: "(4(3(1))(6(5)(7)))",
		},
		{
			desc: "balanced - root",
			tree: New().Insert(4, 2, 6, 1, 3, 5, 7),
			key:  4,
			want: "(5(2(1)(3))(6(7)))",
		},
		{
			desc: "left",
			tree: New().Insert(5, 4, 3, 2, 1),
			key:  3,
			want: "(5(4(2(1))))",
		},
		{
			desc: "unbalanced",
			tree: New().Insert(10, 5, 15, 3, 7, 18),
			key:  10,
			want: "(15(5(3)(7))(18))",
		},
		{
			desc: "hmmm",
			tree: New().Insert(1, 2, 3, 4, 6, 5, 7, 8, 9, 10),
			want: "(1)",
			key:  99,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.tree.Remove(tc.key)
			got := tc.tree.String()
			graphvizPNG(t, tc.tree, fmt.Sprintf("/tmp/tree.png"))
			if got != tc.want {
				t.Errorf("tree mismatch \nwant=%q, \ngot =%q", tc.want, got)
			}

		})
	}
}

// writes the tree as a graphviz file located at fname
func graphvizPNG(t *testing.T, tree *BTree, fname string) {
	dotFilename := fmt.Sprintf("%s.dot", fname)
	f, err := os.OpenFile(dotFilename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	tree.Graphviz(f)
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command("dot", "-Tpng", dotFilename, "-o", fname).Run(); err != nil {
		t.Fatal(err)
	}
}
