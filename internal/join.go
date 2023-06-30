package internal

import (
	"strconv"
	"strings"
)

func JoinInt(sep string, elems ...int) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.Itoa(elems[0])
	}

	var b strings.Builder
	b.WriteString(strconv.Itoa(elems[0]))
	for _, i := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(strconv.Itoa(i))
	}
	return b.String()
}

func JoinUint(sep string, elems ...uint) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.FormatUint(uint64(elems[0]), 10)
	}

	var b strings.Builder
	b.WriteString(strconv.FormatUint(uint64(elems[0]), 10))
	for _, i := range elems[1:] {
		b.WriteString(sep)
		strconv.FormatUint(uint64(i), 10)
	}
	return b.String()
}

func JoinInt8(sep string, elems ...int8) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.FormatInt(int64(elems[0]), 10)
	}

	var b strings.Builder
	b.WriteString(strconv.FormatInt(int64(elems[0]), 10))
	for _, i := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(strconv.FormatInt(int64(i), 10))
	}
	return b.String()
}

func JoinUint8(sep string, elems ...uint8) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.FormatUint(uint64(elems[0]), 10)
	}

	var b strings.Builder
	b.WriteString(strconv.FormatUint(uint64(elems[0]), 10))
	for _, i := range elems[1:] {
		b.WriteString(sep)
		strconv.FormatUint(uint64(i), 10)
	}
	return b.String()
}

func JoinInt16(sep string, elems ...int16) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.FormatInt(int64(elems[0]), 10)
	}

	var b strings.Builder
	b.WriteString(strconv.FormatInt(int64(elems[0]), 10))
	for _, i := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(strconv.FormatInt(int64(i), 10))
	}
	return b.String()
}

func JoinUint16(sep string, elems ...uint16) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.FormatUint(uint64(elems[0]), 10)
	}

	var b strings.Builder
	b.WriteString(strconv.FormatUint(uint64(elems[0]), 10))
	for _, i := range elems[1:] {
		b.WriteString(sep)
		strconv.FormatUint(uint64(i), 10)
	}
	return b.String()
}

func JoinInt32(sep string, elems ...int32) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.FormatInt(int64(elems[0]), 10)
	}

	var b strings.Builder
	b.WriteString(strconv.FormatInt(int64(elems[0]), 10))
	for _, i := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(strconv.FormatInt(int64(i), 10))
	}
	return b.String()
}

func JoinUint32(sep string, elems ...uint32) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.FormatUint(uint64(elems[0]), 10)
	}

	var b strings.Builder
	b.WriteString(strconv.FormatUint(uint64(elems[0]), 10))
	for _, i := range elems[1:] {
		b.WriteString(sep)
		strconv.FormatUint(uint64(i), 10)
	}
	return b.String()
}

func JoinInt64(sep string, elems ...int64) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.FormatInt(elems[0], 10)
	}

	var b strings.Builder
	b.WriteString(strconv.FormatInt(elems[0], 10))
	for _, i := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(strconv.FormatInt(i, 10))
	}
	return b.String()
}

func JoinUint64(sep string, elems ...uint64) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.FormatUint(elems[0], 10)
	}

	var b strings.Builder
	b.WriteString(strconv.FormatUint(elems[0], 10))
	for _, i := range elems[1:] {
		b.WriteString(sep)
		strconv.FormatUint(i, 10)
	}
	return b.String()
}
