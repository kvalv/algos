package page

import (
	"fmt"
	"slices"
)

type slot struct {
	p    CellPointer
	size CellSize

	debugValue byte
}

func (s *slot) String() string {
	return fmt.Sprintf("<slot p=%d size=%d>", s.p, s.size)
}

func (s *slot) WithDebugValue(v byte) *slot {
	if s == nil {
		return nil
	}
	s.debugValue = v
	return s
}

// func (s slot)

func (s slot) End() CellPointer {
	return s.p.Add(s.size)
}
func (s slot) AdjacentTo(p CellPointer) bool {
	return s.p.AdjacentTo(p) || s.End().AdjacentTo(p)
}

type FreeSlots struct {
	size    int
	slots   []slot
	yielded []*slot
}

func NewFreeSlots(size int) *FreeSlots {
	return &FreeSlots{
		size: size,
		slots: []slot{
			{p: 0, size: CellSize(size)},
		},
	}
}

func (f *FreeSlots) Reserve(size CellSize) *slot {
	var (
		best  *slot
		index int
	)
	for i, s := range f.slots {
		if s.size >= size {
			if best == nil || s.size < best.size {
				best = &f.slots[i]
				index = i
			}
		}
	}
	var res *slot
	if best != nil {
		fmt.Printf("best is %s\n", best)
		if n := best.size - size; n > 0 {
			res = &slot{p: best.End().Sub(size), size: size}
			best.size -= size
			fmt.Printf("before return \n\t%s \n\t%s\n", res, best)
		} else {
			tmp := *best
			res = &tmp
			// bug here, we want to only give as little as necessary and keep rest
			f.slots = slices.Delete(f.slots, index, index+1)
		}
	}
	// ...
	if res != nil {
		fmt.Printf("adding to yielded %s\n", res)
		f.yielded = append(f.yielded, res)
	}

	return res
}

func (f *FreeSlots) Free(sl *slot) {
	// are we next to adjacent free slots? In that case, we extend them.
	// Otherwise, we just register
	p := sl.p
	size := sl.size

	for i, slot := range f.slots {
		if p > slot.p {
			if slot.AdjacentTo(p) {
				slot.size += size
				if i+1 < len(f.slots)-1 {
					right := f.slots[i+1]
					if right.AdjacentTo(slot.End()) {
						slot.size += right.size
						f.slots = slices.Delete(f.slots, i+1, i+2)
					}
				}
				return
			}
		}
	}
}

// Primarily for testing
func (f *FreeSlots) String() string {
	b := make([]byte, f.size)
	for i := range b {
		b[i] = '.'
	}
	for _, s := range f.slots {
		fmt.Printf("visiting slot %d %d\n", s.p, s.size)
		for j := 0; j < int(s.size); j++ {
			i := int(s.p)
			if s.debugValue != 0 {
				b[i+j] = s.debugValue
			} else {
				b[i+j] = '.'
			}
		}
	}
	for _, s := range f.yielded {
		fmt.Printf("visiting yielded slot %d %d\n", s.p, s.size)
		for j := 0; j < int(s.size); j++ {
			i := int(s.p)
			if s.debugValue != 0 {
				b[i+j] = s.debugValue
			} else {
				b[i+j] = 'x'
			}
		}
	}
	fmt.Printf("b=%q\n", string(b))
	return string(b)
}
