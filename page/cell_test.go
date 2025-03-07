package page

import (
	"testing"
)

func TestBytes(t *testing.T) {
	cases := []struct {
		desc string
		cell Cell
		want []byte
	}{
		{
			desc: "KeyValue",
			cell: NewValueCell("hi", []byte("world")),
			want: []byte{2, 5, 'h', 'i', 'w', 'o', 'r', 'l', 'd'},
		},
		{
			desc: "Key",
			cell: NewKeyCell("hi", 8),
			want: []byte{2, 0x00, 0x08, 'h', 'i'},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.cell.Bytes()
			expectBytesEq(t, tc.want, got)
		})
	}
}
