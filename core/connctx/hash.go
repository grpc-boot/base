package connctx

import "math"

type Hash map[string]interface{}

func (h Hash) del(fields ...string) (delNum int) {
	for _, field := range fields {
		if _, exists := h[field]; exists {
			delNum++
			delete(h, field)
		}
	}
	return
}

func (h Hash) len() int {
	return len(h)
}

func (h Hash) get(field string) (value interface{}, exists bool) {
	value, exists = h[field]
	return
}

func (h Hash) set(field string, value interface{}) (isCreate bool) {
	_, exists := h[field]

	h[field] = value

	return !exists
}

func (h Hash) mset(fv Hash) {
	if len(fv) < 1 {
		return
	}

	for field, value := range fv {
		h[field] = value
	}
}

func (h Hash) getSet(field string, value interface{}) (old interface{}) {
	old, _ = h[field]
	h[field] = value

	return
}

func (h Hash) setnx(field string, value interface{}) (ok bool) {
	_, exists := h[field]
	if exists {
		return
	}

	h[field] = value

	return true
}

func (h Hash) getall() Hash {
	hh := make(Hash, len(h))
	if len(h) < 1 {
		return hh
	}

	for field, val := range h {
		hh[field] = val
	}

	return hh
}

func (h Hash) incrBy(field string, value int64) (newValue int64, err error) {
	oldValue, exists := h[field]
	if !exists {
		newValue = value
		h[field] = value
		return
	}

	switch val := oldValue.(type) {
	case int64:
		newValue = val + value
		h[field] = newValue
	case int:
		newValue = int64(val) + value
		if newValue > math.MaxInt {
			newValue = int64(val)
			err = ErrValueOutOfRange
			return
		}
		h[field] = int(newValue)
	default:
		err = ErrType
	}
	return
}

func (h Hash) setBit(field string, offset uint16, value bool) (oldValue bool, err error) {
	index := getIndex(offset)
	val, exists := h[field]
	if !exists {
		if !value {
			return
		}

		bm := make([]byte, index+1, index+1)
		bm[index] = getByte(offset)
		h[field] = bm
		return
	}

	v, ok := val.([]byte)
	if !ok {
		err = ErrType
		return
	}

	if index >= len(v) {
		if !value {
			return
		}

		nb := make([]byte, index+1, index+1)
		copy(nb, v)
		v = nb
	} else {
		oldValue = (v[index] & getByte(offset)) > 0
	}

	if !value {
		v[index] = v[index] & (^getByte(offset))
		h[field] = v
		return
	}

	v[index] = v[index] | getByte(offset)
	h[field] = v
	return
}

func (h Hash) getBit(field string, offset uint16) (value bool, err error) {
	val, exists := h[field]
	if !exists {
		return
	}

	v, ok := val.([]byte)
	if !ok {
		err = ErrType
		return
	}

	index := getIndex(offset)
	if index >= len(v) {
		return
	}

	value = (v[index] & getByte(offset)) > 0
	return
}

func (h Hash) hasBit(field string) (has bool, err error) {
	value, exists := h[field]
	if !exists {
		return
	}

	val, ok := value.([]byte)
	if !ok {
		err = ErrType
		return
	}

	for index, _ := range val {
		if val[index] > 0 {
			has = true
			break
		}
	}
	return
}

func (h Hash) bitCount(field string) (num int, err error) {
	value, exists := h[field]
	if !exists {
		return
	}

	val, ok := value.([]byte)
	if !ok {
		err = ErrType
		return
	}

	for index, _ := range val {
		if val[index] == 0 {
			continue
		}

		for i := 0; i < 8; i++ {
			if val[index]&(1<<i) > 0 {
				num++
			}
		}
	}

	return
}

func (h Hash) sAdd(field string, items ...interface{}) (newNum int, err error) {
	value, exists := h[field]
	if !exists {
		s := set{}
		newNum = s.add(items...)
		h[field] = s
		return
	}

	if _, ok := value.(set); !ok {
		err = ErrType
		return
	}

	newNum = h[field].(set).add(items...)
	return
}

func (h Hash) sCard(field string) (total int, err error) {
	if h.len() < 1 {
		return
	}

	value, exists := h[field]
	if !exists {
		return
	}

	if v, ok := value.(set); ok {
		total = v.card()
	} else {
		err = ErrType
	}

	return
}

func (h Hash) sMembers(field string) (items []interface{}, err error) {
	value, exists := h[field]
	if !exists {
		return
	}

	if val, ok := value.(set); ok {
		items = val.members()
	} else {
		err = ErrType
	}

	return
}

func (h Hash) sIsMember(field string, item interface{}) (isMem bool, err error) {
	value, exists := h[field]
	if !exists {
		return
	}

	if val, ok := value.(set); ok {
		isMem = val.isMember(item)
	} else {
		err = ErrType
	}

	return
}

func (h Hash) sRem(field string, items ...interface{}) (delNum int, err error) {
	value, exists := h[field]
	if !exists {
		return
	}

	if _, ok := value.(set); !ok {
		err = ErrType
		return
	}

	delNum = value.(set).rem(items...)
	return
}

func (h Hash) sPop(field string) (item interface{}, err error) {
	value, exists := h[field]
	if !exists {
		return
	}

	if _, ok := value.(set); !ok {
		err = ErrType
		return
	}

	item = value.(set).pop()
	return
}

func (h Hash) sRandMember(field string) (item interface{}, err error) {
	value, exists := h[field]
	if !exists {
		return
	}

	if val, ok := value.(set); ok {
		item = val.randMember()
	} else {
		err = ErrType
	}
	return
}
