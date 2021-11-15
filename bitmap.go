package base

import (
	"bytes"
	"fmt"
)

/**
 * 01234567 89...
 * 10010010 10010001
 */

type Bitmap interface {
	HasBit() bool
	BitCount() int
	Exists(tag int) (exists bool)
	AddTag(tag int) Bitmap
	DelTag(tag int) Bitmap
	Data() (data []byte)
	SprinfBinary() string
}

/**
 * 默认实现
 */
type bitmap struct {
	Bitmap

	data []byte
}

func NewBitmap(data []byte) Bitmap {
	return &bitmap{
		data: data,
	}
}

func (b *bitmap) HasBit() bool {
	if len(b.data) < 1 {
		return false
	}

	for _, bt := range b.data {
		if bt > 0 {
			return true
		}
	}

	return false
}

func (b *bitmap) BitCount() int {
	if len(b.data) < 1 {
		return 0
	}

	bitCount := 0
	for _, bt := range b.data {
		for i := 0; i < 8; i++ {
			if bt&(1<<i) > 0 {
				bitCount++
			}
		}
	}

	return bitCount
}

func (b *bitmap) Exists(tag int) (exists bool) {
	index := tag / 8

	if index >= len(b.data) {
		return false
	}

	return (b.data[index] & uint8(1<<(7-(tag%8)))) > 0
}

func (b *bitmap) AddTag(tag int) Bitmap {
	index := tag / 8

	if index >= len(b.data) {
		nb := make([]byte, index+1, index+1)
		copy(nb, b.data)
		b.data = nb
	}

	b.data[index] = b.data[index] | uint8(1<<(7-(tag%8)))
	return b
}

func (b *bitmap) DelTag(tag int) Bitmap {
	index := tag / 8
	if index >= len(b.data) {
		return b
	}

	b.data[index] = b.data[index] & (^uint8(1 << (7 - (tag % 8))))
	return b
}

func (b *bitmap) Data() (data []byte) {
	return b.data
}

func (b *bitmap) SprinfBinary() string {
	if len(b.data) < 1 {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	for _, d := range b.data {
		buf.WriteString(fmt.Sprintf("%08b ", d))
	}

	data := buf.String()
	return data[:len(data)-1]
}
