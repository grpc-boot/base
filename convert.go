package base

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"unicode"

	"github.com/grpc-boot/base/internal"
)

const (
	codeSalt4Six   = 714924305
	codeSalt4Seven = 42180533654
	codeSalt4Eight = 2488651484856
)

var (
	alphanumericLength byte
	alphanumericSet    = []byte("M0Q1EK4aFzPcZeUgAfXi3SjTkYmnGpqrJstD6uLv9xy7bB5CHNR8VhW2wdI")
	alphanumericMap    = make(map[byte]byte, len(alphanumericSet))
)

func init() {
	alphanumericLength = byte(len(alphanumericSet))

	for i, b := range alphanumericSet {
		alphanumericMap[b] = byte(i)
	}
}

// ToString 转为字符串类型
func ToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return Bytes2String(v)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Bytes2String 字节切片转换为字符串
func Bytes2String(data []byte) string {
	return internal.Bytes2String(data)
}

// Bytes2Int64 字节切片转换为int64
func Bytes2Int64(data []byte) int64 {
	return internal.Bytes2Int64(data)
}

// Bytes2Uint32 字节切片转换为uint32
func Bytes2Uint32(data []byte) uint32 {
	return internal.Bytes2Uint32(data)
}

// Bytes2Uint64 字节切片转换为uint64
func Bytes2Uint64(data []byte) uint64 {
	return internal.Bytes2Uint64(data)
}

// Bytes2Float64 字节切片转换为float64
func Bytes2Float64(data []byte) float64 {
	return Bytes2Float64(data)
}

// Int64ToHex _
func Int64ToHex(value int64) string {
	return strconv.FormatInt(value, 16)
}

// Uint64ToHex _
func Uint64ToHex(value uint64) string {
	return strconv.FormatUint(value, 16)
}

// Int64ToHexWithPad _
func Int64ToHexWithPad(value int64, padLength int) string {
	hexStr := Int64ToHex(value)
	if len(hexStr) >= padLength {
		return hexStr
	}

	buffer := make([]byte, padLength, padLength)

	for i := 0; i < (padLength - len(hexStr)); i++ {
		buffer[i] = '0'
	}

	start := padLength - len(hexStr)
	for j := start; j < padLength; j++ {
		buffer[j] = hexStr[j-start]
	}

	return Bytes2String(buffer)
}

func PackUin32(value uint32) []byte {
	buffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(buffer, binary.BigEndian, value)
	return buffer.Bytes()
}

func PackIn32(value int32) []byte {
	buffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(buffer, binary.BigEndian, value)
	return buffer.Bytes()
}

func UnpackUint32(data []byte) (value uint32, err error) {
	buffer := bytes.NewBuffer(data)
	var val uint32
	err = binary.Read(buffer, binary.BigEndian, &val)
	return val, err
}

func UnpackInt32(data []byte) (value int32, err error) {
	buffer := bytes.NewBuffer(data)
	var val int32
	err = binary.Read(buffer, binary.BigEndian, &val)
	return val, err
}

// Hex2Int64 _
func Hex2Int64(value string) int64 {
	val, _ := strconv.ParseInt(value, 16, 64)
	return val
}

// Hex2Uint64 _
func Hex2Uint64(value string) uint64 {
	val, _ := strconv.ParseUint(value, 16, 64)
	return val
}

// BigCamels 转换为大驼峰
func BigCamels(sep byte, data string) string {
	var (
		fd    = []byte(data)
		upper = true
	)

	for index := 0; index < len(fd); index++ {
		if upper {
			fd[index] = byte(unicode.ToUpper(rune(fd[index])))
			upper = false
			continue
		}

		if fd[index] == sep {
			fd = append(fd[0:index], fd[index+1:]...)
			upper = true
			index--
		}
	}

	return Bytes2String(fd)
}

// SmallCamels 转换为小驼峰
func SmallCamels(sep byte, data string) string {
	var (
		fd    = []byte(data)
		upper = false
	)

	fd[0] = byte(unicode.ToLower(rune(fd[0])))

	for index := 1; index < len(fd); index++ {
		if upper {
			fd[index] = byte(unicode.ToUpper(rune(fd[index])))
			upper = false
			continue
		}

		if fd[index] == sep {
			fd = append(fd[0:index], fd[index+1:]...)
			upper = true
			index--
		}
	}

	return Bytes2String(fd)
}

func LcFirst(str string) string {
	return internal.LcFirst(str)
}

func UcFirst(str string) string {
	return internal.UcFirst(str)
}

func Id2Code6(id uint64) (code string, err error) {
	return id2Code(id+codeSalt4Six, 6)
}

func Id2Code7(id uint64) (code string, err error) {
	return id2Code(id+codeSalt4Seven, 7)
}

func Id2Code8(id uint64) (code string, err error) {
	return id2Code(id+codeSalt4Eight, 8)
}

func id2Code(id uint64, length byte) (code string, err error) {
	if id < 1 {
		return "", ErrOutOfRange
	}

	var (
		codeBytes = make([]byte, length)
	)

	max := uint64(math.Pow(float64(len(alphanumericSet)), float64(length))) - 1
	if id > max {
		return "", ErrOutOfRange
	}

	var (
		c      = byte(id % uint64(len(alphanumericSet)))
		i byte = 1
	)

	codeBytes[0] = alphanumericSet[c]
	id = id / uint64(alphanumericLength)

	for ; i < length; i++ {
		a := byte(id % uint64(alphanumericLength))
		d := (a + i + c) % alphanumericLength

		//fmt.Printf("D: %d A:%d I:%d C: %d\n", d, a, i, c)

		codeBytes[i] = alphanumericSet[d]
		id = id / uint64(len(alphanumericSet))
	}

	return Bytes2String(codeBytes), nil
}

func Code2Uint64(code string) (id uint64, err error) {
	var (
		salt       uint64
		codeBytes  = []byte(code)
		codeLength = byte(len(codeBytes))
	)

	switch codeLength {
	case 6:
		salt = codeSalt4Six
	case 7:
		salt = codeSalt4Seven
	case 8:
		salt = codeSalt4Eight
	default:
		err = ErrOutOfRange
		return
	}

	err = ErrOutOfRange

	c, ok := alphanumericMap[codeBytes[0]]
	if !ok {
		return
	}

	id += uint64(c)

	var i byte = 1

	for ; i < codeLength; i++ {
		d, exists := alphanumericMap[codeBytes[i]]
		if !exists {
			id = 0
			return
		}

		var (
			a     = byte(AbsInt8(int8(d) - int8(i+c)))
			retry = 0
		)

		for a >= alphanumericLength || d != (a+i+c)%alphanumericLength {
			retry++

			if alphanumericLength < a {
				a -= alphanumericLength
				continue
			}

			a = alphanumericLength - a

			if retry > 3 {
				id = 0
				return
			}
		}

		//fmt.Printf("D: %d A:%d I:%d C: %d\n", d, a, i, c)

		id += uint64(a) * uint64(math.Pow(float64(alphanumericLength), float64(i)))
	}

	if id <= salt {
		id = 0
		return
	}

	return id - salt, nil
}
