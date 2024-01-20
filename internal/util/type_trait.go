package util

type Integral interface {
	SignedNumber | UnsignedNumber
}

type UnsignedNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type SignedNumber interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type FloatingNumber interface {
	~float32 | ~float64
}
