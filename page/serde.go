package page

type Serde interface {
	Bytes(h Header) []byte
	FromBytes(any, []byte) error
}
