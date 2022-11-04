package base

import (
	"encoding/binary"
	"hash/crc32"
	"math"
)

const (
	ConnectSuccess = 0x0100
	Tick           = 0x0101

	Login        = 0x0200
	LoginSuccess = 0x0201
	LoginFailed  = 0x0202
)

type Protocol interface {
	Pack(pkg *Package) []byte
	ResponseKey() []byte
	Unpack(data []byte) (pkg *Package, err error)
}

type v0 struct{}

func NewV0() (protocol Protocol, err error) {
	return &v0{}, nil
}

func (pt0 *v0) Pack(pkg *Package) []byte {
	return pkg.Pack()
}

func (pt0 *v0) ResponseKey() []byte {
	return nil
}

func (pt0 *v0) Unpack(data []byte) (pkg *Package, err error) {
	pkg = &Package{}
	err = pkg.Unpack(data)
	return
}

type v1 struct {
	aes *Aes
}

func NewV1(aes *Aes, secretData []byte) (protocol Protocol, err error) {
	data, err := aes.CbcDecrypt(secretData)
	if err != nil {
		return nil, err
	}

	if len(data) != 32 {
		return nil, ErrKeyFormat
	}

	transAes, err := NewAesWithBytes(data[:16], data[16:])
	if err != nil {
		return nil, err
	}

	protocol = &v1{aes: transAes}
	return
}

func (pt1 *v1) header(pkg *Package, data []byte) []byte {
	hexStr := Int64ToHexWithPad(int64(pkg.Id), 4)

	sign := crc32.ChecksumIEEE(data)

	seqBuf := make([]byte, 5, 5)
	binary.PutUvarint(seqBuf, uint64(sign))

	header := make([]byte, 0, 11)
	header = append(header, 1, ':')
	header = append(header, hexStr...)
	header = append(header, seqBuf...)
	return header
}

func (pt1 *v1) Pack(pkg *Package) []byte {
	data := pt1.aes.CbcEncrypt(pkg.Pack())
	header := pt1.header(pkg, data)
	return append(header, data...)
}

func (pt1 *v1) ResponseKey() []byte {
	return nil
}

func (pt1 *v1) Unpack(data []byte) (pkg *Package, err error) {
	if len(data) < 1 {
		return nil, ErrDataEmpty
	}

	if len(data) < 12 || data[0] != 1 || data[1] != ':' {
		return nil, ErrDataFormat
	}

	idInt64 := Hex2Int64(Bytes2String(data[2:6]))
	if idInt64 < 1 || idInt64 > math.MaxUint16 {
		return nil, ErrDataFormat
	}

	sign, _ := binary.Uvarint(data[6:11])
	if crc32.ChecksumIEEE(data[11:]) != uint32(sign) {
		return nil, ErrDataSign
	}

	jsonData, err := pt1.aes.CbcDecrypt(data[11:])
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

type v2 struct {
	aes *Aes
	iv  []byte
}

func NewV2ForClient(aes *Aes, key, secretData []byte) (protocol Protocol, err error) {
	data, err := aes.CbcDecrypt(secretData)
	if err != nil {
		return nil, err
	}

	if len(data) != 16 {
		return nil, ErrKeyFormat
	}

	transAes, err := NewAesWithBytes(key, data)
	if err != nil {
		return nil, err
	}

	protocol = &v2{aes: transAes}
	return
}

func NewV2(aes *Aes, secretData []byte) (protocol Protocol, err error) {
	data, err := aes.CbcDecrypt(secretData)
	if err != nil {
		return nil, err
	}

	if len(data) != 16 {
		return nil, ErrKeyFormat
	}

	iv := RandBytes(16)

	transAes, err := NewAesWithBytes(data, iv)
	if err != nil {
		return nil, err
	}

	protocol = &v2{aes: transAes, iv: aes.CbcEncrypt(iv)}
	return
}

func (pt2 *v2) ResponseKey() []byte {
	return pt2.iv
}

func (pt2 *v2) header(pkg *Package, data []byte) []byte {
	hexStr := Int64ToHexWithPad(int64(pkg.Id), 4)

	sign := crc32.ChecksumIEEE(data)

	seqBuf := make([]byte, 5, 5)
	binary.PutUvarint(seqBuf, uint64(sign))

	header := make([]byte, 0, 11)
	header = append(header, 2, ':')
	header = append(header, hexStr...)
	header = append(header, seqBuf...)
	return header
}

func (pt2 *v2) Pack(pkg *Package) []byte {
	data := pt2.aes.CbcEncrypt(pkg.Pack())
	header := pt2.header(pkg, data)
	return append(header, data...)
}

func (pt2 *v2) Unpack(data []byte) (pkg *Package, err error) {
	if len(data) < 1 {
		return nil, ErrDataEmpty
	}

	if len(data) < 12 || data[0] != 2 || data[1] != ':' {
		return nil, ErrDataFormat
	}

	idInt64 := Hex2Int64(Bytes2String(data[2:6]))
	if idInt64 < 1 || idInt64 > math.MaxUint16 {
		return nil, ErrDataFormat
	}

	sign, _ := binary.Uvarint(data[6:11])
	if crc32.ChecksumIEEE(data[11:]) != uint32(sign) {
		return nil, ErrDataSign
	}

	jsonData, err := pt2.aes.CbcDecrypt(data[11:])
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
