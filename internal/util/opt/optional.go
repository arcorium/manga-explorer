package opt

import "time"

type Optional[T any] struct {
	data []T
}

func (o Optional[T]) HasValue() bool {
	return o.data != nil
}

func (o Optional[T]) Value() *T {
	return &o.data[0]
}

func (o Optional[T]) ValueOr(val T) T {
	if o.HasValue() {
		return *o.Value()
	}
	return val
}

var (
	NullStr  = Null[string]()
	NullTime = Null[time.Time]()
)

func Null[T any]() Optional[T] {
	return Optional[T]{nil}
}

func New[T any](val T) Optional[T] {
	return Optional[T]{
		[]T{val},
	}
}
