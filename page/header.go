package page

type Header struct {
	PageSize uint16   // page size
	CType    CellType // homogenous cell type within a page
}

func (p *Header) DiskSize() int {
	return 3
}

func (p *Header) Bytes(h Header) []byte {
	return nil
}
