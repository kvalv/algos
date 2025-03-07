package page

import "encoding/binary"

type CellType int

// The guys in front of the page
type CellPointer uint16

func (c CellPointer) Add(s CellSize) CellPointer { return CellPointer(uint16(c) + uint16(s)) }
func (c CellPointer) Sub(s CellSize) CellPointer { return CellPointer(uint16(c) - uint16(s)) }

func (c CellPointer) AdjacentTo(other CellPointer) bool {
	if uint16(c)+1 == uint16(other) {
		return true
	}
	if uint16(c)-1 == uint16(other) {
		return true
	}
	return false
}

type CellSize uint16

const (
	CellTypeKey CellType = iota
	CellTypeValue
)

// and also pointer, pairs
type Cell struct {
	Type CellType
	Key  string

	Value  []byte
	PageID PageID
}

func NewKeyCell(key string, ID PageID) Cell {
	return Cell{Type: CellTypeKey, Key: key, PageID: ID}
}
func NewValueCell(key string, value []byte) Cell {
	return Cell{Type: CellTypeValue, Key: key, Value: value}
}

func (c *Cell) DiskSize() int {
	if c.Type == CellTypeKey {
		// keylen, pageID, keydata
		return 1 + 2 + len(c.Key)
	}
	// keylen, valuelen, keydata, valuedata
	return 1 + 1 + len(c.Key) + len(c.Value)
}

func (c *Cell) Write(b []byte) (n int, err error) {
	b[n] = uint8(len(c.Key))
	n++
	if c.Type == CellTypeValue {
		b[n] = uint8(len(c.Value))
		n++
	} else {
		binary.BigEndian.PutUint16(b[n:n+2], uint16(c.PageID))
		n++
		n++
	}
	for _, c := range c.Key {
		b[n] = byte(c) // assume ascii
		n++
	}
	if c.Type == CellTypeValue {
		for _, c := range c.Value {
			b[n] = byte(c)
			n++
		}
	}
	return n, err
}

// utility function, for testing
func (c *Cell) Bytes() []byte {
	b := make([]byte, c.DiskSize())
	n, err := c.Write(b)
	if n != len(b) {
		panic("Bytes: byte length mismatch")
	}
	if err != nil {
		panic(err)
	}
	return b
}
