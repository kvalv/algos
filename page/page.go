package page

// cell -> page -> tree
// one tree on many pages; many cells into one page

type PageID int

type Header struct {
	CType CellType
}

type Page struct {
	Header  Header
	Offsets []CellPointer
}

func (p *Page) Bytes(h Header) {
	// ..
}
