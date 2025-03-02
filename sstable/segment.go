package sstable

import (
	"slices"
)

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
	var res Segment
	for _, rec := range s.records {
		if i := slices.IndexFunc(res.records, func(r Record) bool {
			return r.Key == rec.Key
		}); i == -1 {
			res.records = append(res.records, rec)
		} else {
			res.records[i] = rec
		}
	}
	return &res
}
