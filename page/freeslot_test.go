package page

import "testing"

func TestFreeSlot(t *testing.T) {
	cases := []struct {
		desc string
		size int
		init func(fs *FreeSlots)
		want any
	}{
		{
			desc: "simple",
			size: 5,
			init: func(fs *FreeSlots) {
				fs.Reserve(3).WithDebugValue('a')
			},
			want: "..aaa",
		},
		{
			desc: "two reserves",
			size: 5,
			init: func(fs *FreeSlots) {
				fs.Reserve(3).WithDebugValue('a')
				fs.Reserve(1).WithDebugValue('b')
			},
			want: ".baaa",
		},
		{
			desc: "reserve and free",
			size: 5,
			init: func(fs *FreeSlots) {
				a := fs.Reserve(3).WithDebugValue('a')
				fs.Reserve(1).WithDebugValue('b')
				fs.Free(a)
				fs.Reserve(2).WithDebugValue('c')
			},
			want: ".b.cc",
		},
		{
			desc: "reserve and fragment",
			size: 5,
			init: func(fs *FreeSlots) {
				a := fs.Reserve(1).WithDebugValue('a')
				b := fs.Reserve(1).WithDebugValue('b')
				fs.Free(a)
				fs.Free(b)
				fs.Reserve(2).WithDebugValue('c')
			},
			want: "...cc",
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			fs := NewFreeSlots(tc.size)
			tc.init(fs)
			got := fs.String()
			if got != tc.want {
				t.Fatalf("want=%q, got=%q", tc.want, got)
			}
		})
	}

}
