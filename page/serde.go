package page

import "io"

type Serde interface {
	io.Writer
	// Bytes(h Header) []byte
	FromBytes(any, []byte) error
}
