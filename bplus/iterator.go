package bplus

type Iterator[T any] interface {
	Next() *T
}

type RangeIterator = Iterator[Match]

type iterator[T any] struct {
	next func() *T
}

func (i *iterator[T]) Next() *T {
	return i.next()
}

func NewIterator[T any](next func() *T) Iterator[T] {
	return &iterator[T]{next: next}
}
func NewEmptyIterator[T any]() Iterator[T] {
	return NewIterator(func() *T { return nil })
}
