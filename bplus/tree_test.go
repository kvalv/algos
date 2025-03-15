package bplus

import (
	"fmt"
	"os"
	"testing"
)

func TestFromString(t *testing.T) {
	tree := FromString(2, "(2(1)(34))", os.Stderr)
	want := "(2(1)(34))"
	expectTree(t, want, tree)
}

func TestFind(t *testing.T) {
	cases := []struct {
		input string
		key   int
		want  string
		index int
	}{
		{
			input: "(2(1)(34))", key: 4,
			want: "(34)", index: 1,
		},
		{
			input: "(58(23)(67)(9))", key: 6,
			want: "(67)", index: 0,
		},
		{
			input: "(5(12)(78))", key: 99,
			want: "",
		},
		{
			input: "(9(12)(89))", key: 9,
			want: "(89)", index: 1,
		},
		{
			input: "(5(1)(8))", key: 1,
			want: "(1)", index: 0,
		},
		{
			input: "(5(1)(8))", key: 8,
			want: "(8)", index: 0,
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s/%d", tc.input, tc.key), func(t *testing.T) {
			tree := FromString(2, tc.input, os.Stderr)
			got := tree.Find(tc.key)
			var want *Match
			if tc.want != "" {
				want = &Match{
					Node:  FromString(2, tc.want, os.Stderr).Root,
					Index: tc.index,
				}
			}
			expectMatch(t, want, got)
		})
	}
}

type testMatch struct {
	nodestr string
	index   int
}

func (m *testMatch) String() string {
	if m == nil {
		return "nil"
	}
	return fmt.Sprintf("%s/%d", m.nodestr, m.index)
}

func TestRange(t *testing.T) {
	cases := []struct {
		input        string
		lower, upper int
		want         []string
	}{
		{
			input: "(3(12)(34))",
			lower: 2, upper: 4,
			want: []string{"(12)/1", "(34)/0"},
		},
		{
			input: "(3(12)(34))", // no nothin'
			lower: 5, upper: 12,
			want: nil,
		},
		{
			input: "(3(12)(34))", // full
			lower: -1, upper: 12,
			want: []string{"(12)/0", "(12)/1", "(34)/0", "(34)/1"},
		},
		{
			input: "(123)", // root only
			lower: 1, upper: 3,
			want: []string{"(123)/0", "(123)/1"},
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s/%d/%d", tc.input, tc.lower, tc.upper), func(t *testing.T) {
			tree := FromString(2, tc.input, os.Stderr)
			got := tree.Range(tc.lower, tc.upper)
			expectMatches(t, tc.want, got)
		})
	}
}

func TestInsert(t *testing.T) {
	cases := []struct {
		input string
		key   int
		want  string
	}{
		{
			input: "(ab)",
			key:   'c',
			want:  "(abc)",
		},
		{
			input: "(bcd)",
			key:   'a',
			want:  "(c(ab)(cd))",
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s/%d", tc.input, tc.key), func(t *testing.T) {
			tree := FromString(2, tc.input, os.Stderr)
			tree.Insert(tc.key, 0)
			expectTree(t, tc.want, tree)
		})
	}
}

func expectMatches(t *testing.T, want []string, got Iterator[Match]) {
	t.Helper()
	for i, w := range want {
		gotMatch := got.Next()
		if gotMatch == nil {
			t.Fatalf("missing match; want=%s", w)
		}
		if s := gotMatch.String(); s != w {
			t.Fatalf("mismatch at index %d: want=%s, got=%s", i, w, s)
		}
	}
	if val := got.Next(); val != nil {
		t.Fatalf("unexpected extra match; got = %v", val)
	}
}

func expectMatch(t *testing.T, want, got *Match) {
	t.Helper()
	if want == nil {
		if got != nil {
			t.Fatalf("want nil, got non-nil node %v", got)
		}
		return
	}
	if got.Node.String() != want.Node.String() {
		t.Fatalf("unexpected node structure;\nwant= %s\ngot = %s", want.Node.String(), got.Node.String())
	}
	if got.Index != want.Index {
		t.Fatalf("unexpected index; want = %d, got = %d", want.Index, got.Index)
	}
}

func expectTree(t *testing.T, want string, got *BTree) {
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
