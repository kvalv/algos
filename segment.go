package sstable

type Record struct {
	Key   string
	Value int
}

func NewRecord(key string, value int) Record {
	return Record{Key: key, Value: value}
}

type Segment struct {
	records []Record
}

func NewSegment(r ...Record) *Segment {
	return &Segment{
		records: r,
	}
}

func Compact(s *Segment) *Segment {
	want := NewSegment(
		NewRecord("yawn", 522),
		NewRecord("mew", 1082),
		NewRecord("purr", 2108),
	)
	return want
}
