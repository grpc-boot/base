package base

import (
	"encoding/base64"
	"encoding/hex"
)

const (
	LevelJson = 0
	LevelV1   = 1
	LevelV2   = 2
)

type Accept struct {
	aes     *Aes
	level   uint8
	protoV0 Protocol
}

func NewAccept(aes *Aes, level uint8) *Accept {
	v0Proto, _ := NewV0()

	return &Accept{
		aes:     aes,
		level:   level,
		protoV0: v0Proto,
	}
}

func (a *Accept) Accept(level uint8, secretData []byte) (protocol Protocol, err error) {
	if level < a.level || level < LevelJson || level > LevelV2 {
		return nil, ErrForbidden
	}

	switch level {
	case LevelV1:
		return NewV1(a.aes, secretData)
	case LevelV2:
		return NewV2(a.aes, secretData)
	default:
		return a.protoV0, nil
	}
}

func (a *Accept) AcceptHex(level uint8, hexData []byte) (protocol Protocol, err error) {
	var data []byte
	if len(hexData) > 0 {
		data, err = hex.DecodeString(Bytes2String(hexData))
		if err != nil {
			return nil, ErrForbidden
		}
	}
	return a.Accept(level, data)
}

func (a *Accept) AcceptBase64Url(level uint8, base64Data []byte, urlSafe bool) (protocol Protocol, err error) {
	var data []byte
	if len(base64Data) > 0 {
		if urlSafe {
			data, err = base64.URLEncoding.DecodeString(Bytes2String(base64Data))
		} else {
			data, err = base64.StdEncoding.DecodeString(Bytes2String(base64Data))
		}

		if err != nil {
			return nil, ErrForbidden
		}
	}

	return a.Accept(level, data)
}
