package page

import "testing"

func TestInsert(t *testing.T) {
	p, _ := NewPage(30)

	if _, err := p.Insert(NewKeyCell("foo", 10)); err != nil {
		t.Fatalf("failed to insert: %s", err)
	}

	if _, err := p.Insert(NewValueCell("bar", []byte("xx"))); err != nil {
		t.Fatalf("failed to insert: %s", err)
	}

	got := make([]byte, p.Header.PageSize)
	if _, err := p.Write(got); err != nil {
		t.Fatalf("failed to write: %s", err)
	}

	want := []byte{
		30, 0x01, // header
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 30, // pointer to first cell
		0, 10, // pointer to 2nd cell
		3, 10, 'f', 'o', 'o', // first cell
		3, 2, 'b', 'a', 'r', 'x', 'x', // 2nd cell
	}
	if len(want) != 30 {
		t.Fatalf("expected 30, got %d", len(want)) // sanity check
	}
	expectBytesEq(t, want, got)
}
