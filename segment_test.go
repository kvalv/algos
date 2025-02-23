package sstable

import (
	"testing"
)

var rec = NewRecord

func TestSegment(t *testing.T) {
	seg := NewSegment(
		rec("mew", 1078),
		rec("purr", 2103),
		rec("purr", 2104),
		rec("mew", 1079),
		rec("mew", 1080),
		rec("mew", 1081),
		rec("purr", 2105),
		rec("purr", 2106),
		rec("purr", 2107),
		rec("yawn", 522),
		rec("purr", 2108),
		rec("mew", 1082),
	)
	got := Compact(seg)
	want := NewSegment(
		rec("yawn", 522),
		rec("mew", 1082),
		rec("purr", 2108),
	)
	expectSegment(t, want, got)
}

func expectSegment(t *testing.T, want, got *Segment) {
	t.Helper()
	if got == nil {
		t.Fatalf("segment is nil")
	}
	if len(want.records) != len(got.records) {
		t.Errorf("got %d records, want %d", len(got.records), len(want.records))
	}
	for i, w := range want.records {
		g := got.records[i]
		if w.Key != g.Key {
			t.Errorf("key mismatch; want=%q, got=%q", w.Key, g.Key)
		}
		if w.Value != g.Value {
			t.Errorf("value mismatch; want=%d, got=%d", w.Value, g.Value)
		}
	}
}
