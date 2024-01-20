package util

type Enum[T any] interface {
	String() string
	Underlying() T
	Validate() error
}
