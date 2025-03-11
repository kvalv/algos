package bplus

import (
	"fmt"
	"os"
	"testing"
)

func TestFromString(t *testing.T) {
	tree := FromString(2, "(2(1)(34))", os.Stderr)
	want := "(2(1)(34))"
	expectTree(t, tree, want)
}

func TestFind(t *testing.T) {
	tree := FromString(2, "(2(1)(34))", os.Stderr)
	cases := []struct {
		input string
		key   int
		want  *struct {
			index int
			leaf  string
		}
	}{
		{
			input: "(2(1)(34))",
			key:   4,
			want: &struct {
				index int
				leaf  string
			}{
				index: 1, leaf: "(34)",
			},
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s/%d", tc.input, tc.key), func(t *testing.T) {
			leaf, index := tree.Find(tc.key)
			if tc.want == nil {
				if leaf != nil {
					t.Fatalf("want nil, got non-nil node %s", leaf)
				}
				return
			}
			expectNode(t, leaf, tc.want.leaf)
			if index != tc.want.index {
				t.Fatalf("unexpected index; want = %d, got = %d", tc.want.index, index)
			}
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
func expectNode(t *testing.T, got *Node, want string) {
	t.Helper()
	if want == "" {
		if got != nil {
			t.Fatalf("want nil, got non-nil node %s", got)
		}
		return
	}
	if got.String() != want {
		t.Fatalf("unexpected node structure;\nwant= %s\ngot = %s", want, got.String())
	}
}
