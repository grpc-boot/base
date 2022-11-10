package base

import (
	"encoding/base64"
	"math"
)

const (
	EventConnectSuccess = 0x0100
	EventTick           = 0x0101
	EventClose          = 0x0102
	EventError          = 0x0103
	EventLogin          = 0x0200
	EventLoginSuccess   = 0x0201
	EventLoginFailed    = 0x0202
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
	aes          *Aes
	headerLength int
	startTag     string
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

	protocol = &v1{
		aes:          transAes,
		headerLength: 16,
		startTag:     "013a",
	}
	return
}

func (pt1 *v1) header(pkg *Package, data []byte) []byte {
	hexStr := Int64ToHexWithPad(int64(pkg.Id), 4)
	header := make([]byte, 0, len(data)+pt1.headerLength)
	header = append(header, pt1.startTag...)
	header = append(header, hexStr...)
	header = append(header, Int64ToHexWithPad(int64(len(data)), 8)...)
	return header
}

func (pt1 *v1) Pack(pkg *Package) []byte {
	data := pt1.aes.CbcEncrypt(pkg.Pack())

	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dst, data)
	header := pt1.header(pkg, dst)
	return append(header, dst...)
}

func (pt1 *v1) ResponseKey() []byte {
	return nil
}

func (pt1 *v1) Unpack(data []byte) (pkg *Package, err error) {
	if len(data) < 1 {
		return nil, ErrDataEmpty
	}

	if len(data) < (pt1.headerLength+1) || Bytes2String(data[:4]) != pt1.startTag {
		return nil, ErrDataFormat
	}

	idInt64 := Hex2Int64(Bytes2String(data[4:8]))
	if idInt64 < 1 || idInt64 > math.MaxUint16 {
		return nil, ErrDataFormat
	}

	bodyLength := int(Hex2Int64(Bytes2String(data[8:pt1.headerLength])))
	if bodyLength+pt1.headerLength != len(data) {
		return nil, ErrDataSign
	}

	binaryData, err := base64.StdEncoding.DecodeString(Bytes2String(data[pt1.headerLength:]))
	if err != nil {
		return nil, ErrDataFormat
	}

	jsonData, err := pt1.aes.CbcDecrypt(binaryData)
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
	aes          *Aes
	iv           []byte
	headerLength int
	startTag     string
}

func newV2(aes *Aes) *v2 {
	return &v2{
		aes:          aes,
		headerLength: 16,
		startTag:     "023a",
	}
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

	return newV2(transAes), nil
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

	pt := newV2(transAes)
	pt.iv = aes.CbcEncrypt(iv)

	return pt, nil
}

func (pt2 *v2) ResponseKey() []byte {
	return pt2.iv
}

func (pt2 *v2) header(pkg *Package, data []byte) []byte {
	hexStr := Int64ToHexWithPad(int64(pkg.Id), 4)
	header := make([]byte, 0, len(data)+pt2.headerLength)
	header = append(header, pt2.startTag...)
	header = append(header, hexStr...)
	header = append(header, Int64ToHexWithPad(int64(len(data)), 8)...)
	return header
}

func (pt2 *v2) Pack(pkg *Package) []byte {
	data := pt2.aes.CbcEncrypt(pkg.Pack())

	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dst, data)

	header := pt2.header(pkg, dst)
	return append(header, dst...)
}

func (pt2 *v2) Unpack(data []byte) (pkg *Package, err error) {
	if len(data) < 1 {
		return nil, ErrDataEmpty
	}

	if len(data) < pt2.headerLength+1 || Bytes2String(data[:4]) != pt2.startTag {
		return nil, ErrDataFormat
	}

	idInt64 := Hex2Int64(Bytes2String(data[4:8]))
	if idInt64 < 1 || idInt64 > math.MaxUint16 {
		return nil, ErrDataFormat
	}

	bodyLength := int(Hex2Int64(Bytes2String(data[8:pt2.headerLength])))
	if bodyLength+pt2.headerLength != len(data) {
		return nil, ErrDataSign
	}

	binaryData, err := base64.StdEncoding.DecodeString(Bytes2String(data[pt2.headerLength:]))
	if err != nil {
		return nil, ErrDataFormat
	}

	jsonData, err := pt2.aes.CbcDecrypt(binaryData)
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
