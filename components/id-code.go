package components

import (
	"bytes"
	"math"

	"github.com/grpc-boot/base/v3/utils"
)

const (
	alphanumericMin = 50
	alphanumericMax = 62
)

const (
	defaultSalt6 = 714924305
	defaultSalt8 = 2488651484856
)

var (
	DefaultAlphanumeric = []byte("M0Q1EK4aFzPcZeUgAfXi3SjTkYmnGpqrJstD6uLv9xy7bB5CHNR8VhW2wdI")
	DefaultIdCode, _    = NewIdCode(DefaultAlphanumeric, defaultSalt6, defaultSalt8)
)

func Max6() int64 {
	return DefaultIdCode.Max6()
}

func Max8() int64 {
	return DefaultIdCode.Max8()
}

func Code6(id int64) (code []byte, err error) {
	return DefaultIdCode.Code6(id)
}

func Code6String(id int64) (code string, err error) {
	return DefaultIdCode.Code6String(id)
}

func Code8(id int64) (code []byte, err error) {
	return DefaultIdCode.Code8(id)
}

func Code8String(id int64) (code string, err error) {
	return DefaultIdCode.Code8String(id)
}

func Code2Id(code []byte) (id int64, err error) {
	return DefaultIdCode.Code2Id(code)
}

func CodeString2Id(code string) (id int64, err error) {
	return DefaultIdCode.CodeString2Id(code)
}

func Code6To8(code6 []byte) (code8 []byte, err error) {
	return DefaultIdCode.Code6To8(code6)
}

func CodeString6To8(code6 string) (code8 string, err error) {
	return DefaultIdCode.CodeString6To8(code6)
}

type IdCode struct {
	alphanumericSet    []byte
	alphanumericMap    map[byte]byte
	alphanumericLength byte
	salt6              int64
	salt8              int64
	max6               int64
	max8               int64
}

func NewIdCode(alphanumericSet []byte, salt6, salt8 int64) (*IdCode, error) {
	if salt6 < 1 {
		salt6 = defaultSalt6
	}

	if salt8 < 1 {
		salt8 = defaultSalt8
	}

	if len(alphanumericSet) == 0 {
		alphanumericSet = DefaultAlphanumeric
	}

	ic := &IdCode{
		alphanumericSet: alphanumericSet,
		salt6:           salt6,
		salt8:           salt8,
	}

	err := ic.initAndCheckAlphanumeric()
	if err != nil {
		return nil, err
	}

	return ic, nil
}

func (ic *IdCode) initAndCheckAlphanumeric() error {
	if len(ic.alphanumericSet) < alphanumericMin {
		return ErrAlphanumericLength
	}

	if len(ic.alphanumericSet) > alphanumericMax {
		return ErrAlphanumericLength
	}

	ic.alphanumericLength = byte(len(ic.alphanumericSet))
	ic.alphanumericMap = make(map[byte]byte, ic.alphanumericLength)

	for index, b := range ic.alphanumericSet {
		if _, exists := ic.alphanumericMap[b]; exists {
			return ErrAlphanumeric
		}

		ic.alphanumericMap[b] = byte(index)
	}

	dataBytes := bytes.Repeat([]byte{ic.alphanumericSet[ic.alphanumericLength-1]}, 8)
	ic.max8, _ = ic.Code2Id(dataBytes)
	ic.max6, _ = ic.Code2Id(dataBytes[:6])

	return nil
}

func (ic *IdCode) Max6() int64 {
	return ic.max6
}

func (ic *IdCode) Max8() int64 {
	return ic.max8
}

func (ic *IdCode) Code6String(id int64) (code string, err error) {
	cb, err := ic.Code6(id)
	if err != nil {
		return
	}

	return utils.Bytes2String(cb), err
}

func (ic *IdCode) Code6(id int64) (code []byte, err error) {
	return ic.id2Code(id+ic.salt6, 6)
}

func (ic *IdCode) Code8(id int64) (code []byte, err error) {
	return ic.id2Code(id+ic.salt8, 8)
}

func (ic *IdCode) Code8String(id int64) (code string, err error) {
	cb, err := ic.Code8(id)
	if err != nil {
		return
	}

	return utils.Bytes2String(cb), err
}

func (ic *IdCode) id2Code(id int64, length byte) (code []byte, err error) {
	if id < 1 {
		return nil, ErrOutOfRange
	}

	var (
		codeBytes = make([]byte, length)
	)

	var (
		c      = byte(id % int64(ic.alphanumericLength))
		i byte = 1
	)

	codeBytes[0] = ic.alphanumericSet[c]
	id = id / int64(ic.alphanumericLength)

	for ; i < length; i++ {
		a := byte(id % int64(ic.alphanumericLength))
		d := (a + i + c) % ic.alphanumericLength

		codeBytes[i] = ic.alphanumericSet[d]
		id = id / int64(ic.alphanumericLength)
	}

	return codeBytes, nil
}

func (ic *IdCode) CodeString2Id(code string) (id int64, err error) {
	return ic.Code2Id([]byte(code))
}

func (ic *IdCode) Code2Id(code []byte) (id int64, err error) {
	var (
		salt       int64
		codeBytes  = code
		codeLength = byte(len(codeBytes))
	)

	switch codeLength {
	case 6:
		salt = ic.salt6
	case 8:
		salt = ic.salt8
	default:
		err = ErrOutOfRange
		return
	}

	err = ErrOutOfRange

	c, ok := ic.alphanumericMap[codeBytes[0]]
	if !ok {
		return
	}

	id += int64(c)

	var i byte = 1

	for ; i < codeLength; i++ {
		d, exists := ic.alphanumericMap[codeBytes[i]]
		if !exists {
			id = 0
			return
		}

		var (
			a     = byte(utils.Abs(int8(d) - int8(i+c)))
			retry = 0
		)

		for a >= ic.alphanumericLength || d != (a+i+c)%ic.alphanumericLength {
			retry++

			if ic.alphanumericLength < a {
				a -= ic.alphanumericLength
				continue
			}

			a = ic.alphanumericLength - a

			if retry > 3 {
				id = 0
				return
			}
		}

		id += int64(a) * int64(math.Pow(float64(ic.alphanumericLength), float64(i)))
	}

	if id <= salt {
		id = 0
		return
	}

	return id - salt, nil
}

func (ic *IdCode) Code6To8(code6 []byte) (code8 []byte, err error) {
	id, err := ic.Code2Id(code6)
	if err != nil {
		return
	}

	return ic.Code8(id)
}

func (ic *IdCode) CodeString6To8(code6 string) (code8 string, err error) {
	code, err := ic.Code6To8([]byte(code6))
	if err != nil {
		return
	}

	return utils.Bytes2String(code), err
}
