package page

import (
	"testing"
)

func TestBytes(t *testing.T) {

	cases := []struct {
		desc string
		cell Cell
		head Header
		want []byte
	}{
		{
			desc: "KeyValue",
			cell: Cell{Key: "hi", Value: []byte("world")},
			head: Header{CType: CTKeyValue},
			want: []byte{2, 5, 'h', 'i', 'w', 'o', 'r', 'l', 'd'},
		},
		{
			desc: "Key",
			cell: Cell{Key: "hi", Value: []byte("world")},
			head: Header{CType: CTKey},
			want: []byte{2, 'h', 'i'},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.cell.Bytes(tc.head)
			expectBytesEq(t, tc.want, got)
		})
	}

}

func expectBytesEq(t *testing.T, want, got []byte) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("length mismatch: want=%d, got=%d\n%q", len(want), len(got), got)
	}
	for i := range want {
		if want[i] != got[i] {
			t.Fatalf("mismatch at index %d: want=%v, got=%v (want=%v, got=%v)", i, want[i], got[i], want, got)
		}
	}
}
