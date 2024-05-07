package obj

type Result[T any] struct {
  data T
  err  error
}

func (r Result[T]) Unwrap() (T, error) {
  return r.data, r.err
}

func (r Result[T]) Err() error {
  return r.err
}

func (r Result[T]) Data() T {
  return r.data
}

func Wrap[T any](data T, err error) Result[T] {
  return Result[T]{data: data, err: err}
}
