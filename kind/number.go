package kind

type Int interface {
	~int8 | ~int16 | ~int32 | ~int | ~int64
}

type Uint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint | ~uint64
}

type Integer interface {
	Int | Uint
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	Int | Uint | Float
}
