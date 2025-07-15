package components

import (
	"bytes"
)

var (
	DefaultBase62Charset = []byte("M01EK4aFzPcHNwZeUgAfXi3SjTkYmnGpqrJsuR8VohW2LvQ9xy7bB5CtD6dI")
	DefaultBase62, _     = NewBase62(DefaultBase62Charset)
)

type Base62 interface {
	Encode(value uint64, length int) []byte
	Decode(data []byte) (uint64, error)
}

type base62 struct {
	charset []byte
	charmap [62]int
}

func NewBase62(charset []byte) (Base62, error) {
	if len(charset) > 62 || len(charset) < 50 {
		return nil, ErrAlphanumericLength
	}

	b62 := &base62{
		charset: charset,
	}

	for i, c := range charset {
		index := b62.charIndex(int(c))
		if index >= 63 {
			return nil, ErrAlphanumeric
		}

		b62.charmap[index] = i
	}

	return b62, nil
}

func (b *base62) charIndex(c int) int {
	if c >= 'A' && c <= 'Z' {
		return c - 65 + 10
	}

	if c >= 'a' && c <= 'z' {
		return c - 97 + 10 + 26
	}

	if c >= '0' && c <= '9' {
		return c - 48
	}

	return 63
}

func (b *base62) Encode(value uint64, length int) []byte {
	if value == 0 {
		return bytes.Repeat([]byte{b.charset[0]}, length)
	}

	result := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		result[i] = b.charset[value%uint64(len(b.charset))]
		value /= uint64(len(b.charset))
	}

	return result
}

func (b *base62) Decode(code []byte) (uint64, error) {
	var result uint64
	for i := 0; i < len(code); i++ {
		index := b.charIndex(int(code[i]))
		if index >= 63 {
			return 0, ErrDataFormat
		}

		result = result*uint64(len(b.charset)) + uint64(b.charmap[index])
	}

	return result, nil
}
