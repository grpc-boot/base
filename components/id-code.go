package components

import (
	"github.com/grpc-boot/base/v3/utils"
)

const (
	defaultSalt = 714924305
)

var (
	DefaultIdCode, _ = NewIdCode(DefaultBase62, defaultSalt)
)

type IdCode struct {
	base62 Base62
	salt   int64
}

func NewIdCode(base62 Base62, salt int64) (*IdCode, error) {
	ic := &IdCode{
		base62: base62,
		salt:   salt,
	}

	return ic, nil
}

func NewIdCodeWithCharset(charset []byte, salt int64) (*IdCode, error) {
	b62, err := NewBase62(charset)
	if err != nil {
		return nil, err
	}

	return NewIdCode(b62, salt)
}

func (ic *IdCode) Id2Code64(id uint64) (code []byte, err error) {
	if id < 1 {
		return nil, ErrOutOfRange
	}

	obfuscated := utils.FeistelEncrypt64(id, uint64(ic.salt))
	return ic.base62.Encode(obfuscated, 12), nil
}

func (ic *IdCode) Id2Code32(id uint32) (code []byte, err error) {
	if id < 1 {
		return nil, ErrOutOfRange
	}

	obfuscated := utils.FeistelEncrypt32(id, uint32(ic.salt))
	return ic.base62.Encode(uint64(obfuscated), 6), nil
}

func (ic *IdCode) CodeString2Id(code string) (id uint64, err error) {
	return ic.Code2Id([]byte(code))
}

func (ic *IdCode) Code2Id(code []byte) (id uint64, err error) {
	num, err := ic.base62.Decode(code)
	if err != nil {
		return 0, err
	}

	if len(code) == 12 {
		return utils.FeistelDecrypt64(num, uint64(ic.salt)), nil
	}

	return uint64(utils.FeistelDecrypt32(uint32(num), uint32(ic.salt))), nil
}
