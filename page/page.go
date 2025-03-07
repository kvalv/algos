package page

import "fmt"

// cell -> page -> tree
// one tree on many pages; many cells into one page

type PageID uint16

type Page struct {
	Header  Header
	Offsets []CellPointer
	Cells   []Cell
}

func NewPage(size int) (*Page, error) {
	if int(uint16(size)) != size {
		return nil, fmt.Errorf("invalid page size")
	}
	return &Page{
		Header: Header{
			PageSize: uint16(size),
		},
	}, nil
}

func (p *Page) Write(b []byte) (int, error) {
	// b := bytes.Buffer{}
	// io.Writer
	// b.Write(p.Header.Bytes())
	return 0, fmt.Errorf("not implemented")
}

func (p *Page) Insert(cell Cell) (CellPointer, error) {
	return 0, ErrNoSpace
}

func (p *Page) FreeSpace() int {
	n := int(p.Header.PageSize) - p.Header.DiskSize() - len(p.Offsets)*2
	for _, c := range p.Cells {
		n -= c.DiskSize()
	}
	return n
}

func (p *Page) Validate() {
	if len(p.Offsets) != len(p.Cells) {
		panic("offsets and cells mismatch")
	}
	for i, c := range p.Cells {
		if p.Header.CType != c.Type {
			panic(fmt.Sprintf("cell %d has type %v, expected %v", i, c.Type, p.Header.CType))
		}
	}
}
