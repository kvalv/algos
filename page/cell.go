package page

type CellType int

// The guys in front of the page
type CellPointer uint16

const (
	CTKey CellType = iota
	CTKeyValue
)

// and also pointer, pairs
type Cell struct {
	Key   string
	Value []byte
}

// klen, vlen, key, value
func (c *Cell) Bytes(h Header) []byte {
	n := 1 + len(c.Key)
	if h.CType == CTKeyValue {
		n += 1 // vsize
		n += len(c.Value)
	}

	var i int
	b := make([]byte, n)
	b[0] = uint8(len(c.Key))
	i++
	if h.CType == CTKeyValue {
		b[1] = uint8(len(c.Value))
		i++
	}
	for _, c := range c.Key {
		b[i] = byte(c) // assume ascii
		i++
	}
	if h.CType == CTKeyValue {
		for _, c := range c.Value {
			b[i] = byte(c)
			i++
		}
	}
	return b
}
