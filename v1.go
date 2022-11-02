package base

import (
	"encoding/binary"
	"math"

	"go.uber.org/atomic"
)

type V1 struct {
	aes *Aes
	seq atomic.Uint32
}

func NewV1(key string, secretData []byte) (protocol *V1, err error) {
	if len(key) != 32 {
		return nil, ErrKeyFormat
	}

	aes, err := NewAes(key[:16], key[16:])
	if err != nil {
		return nil, err
	}

	data, err := aes.CbcDecrypt(secretData)
	if err != nil {
		return nil, err
	}

	if len(data) != 32 {
		return nil, ErrKeyFormat
	}

	aes, err = NewAesWithBytes(data[:16], data[16:])
	if err != nil {
		return nil, err
	}

	protocol = &V1{aes: aes}
	return
}

func (v1 *V1) header(pkg *Package) []byte {
	hexStr := Int64ToHexWithPad(int64(pkg.Id), 4)

	seq := v1.seq.Inc()
	seqBuf := make([]byte, 4, 4)
	binary.PutUvarint(seqBuf, uint64(seq))

	header := make([]byte, 0, 10)
	header = append(header, 1, ':')
	header = append(header, hexStr...)
	header = append(header, seqBuf...)
	return header
}

func (v1 *V1) Pack(pkg *Package) []byte {
	header := v1.header(pkg)
	return append(header, v1.aes.CbcEncrypt(pkg.Pack())...)
}

func (v1 *V1) Unpack(data []byte) (pkg *Package, err error) {
	if len(data) < 1 {
		return nil, ErrDataEmpty
	}

	if len(data) < 11 || data[0] != 1 || data[1] != ':' {
		return nil, ErrDataFormat
	}

	idInt64 := Hex2Int64(Bytes2String(data[2:6]))
	if idInt64 < 1 || idInt64 > math.MaxUint16 {
		return nil, ErrDataFormat
	}

	jsonData, err := v1.aes.CbcDecrypt(data[10:])
	if err != nil {
		return nil, err
	}

	pkg = &Package{}
	if err = pkg.Unpack(jsonData); err != nil {
		return nil, err
	}

	if pkg.Id != uint16(idInt64) {
		return nil, ErrDataSign
	}
	return
}
