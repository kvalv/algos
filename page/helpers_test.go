package page

import "testing"

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
