package connctx

import (
	"math/rand"
	"testing"
)

func TestCtx_IncrBy(t *testing.T) {
	c := AcquireCtx()
	defer c.Close()
	key := "incr_test"

	val, err := c.IncrBy(key, 100)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	if val != 100 {
		t.Fatalf("want 100, got %d", val)
	}

	val, err = c.Incr(key)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	if val != 101 {
		t.Fatalf("want 101, got %d", val)
	}

	val, _ = c.Decr(key)
	if val != 100 {
		t.Fatalf("want 100, got %d", val)
	}

	val, _ = c.DecrBy(key, 100)
	if val != 0 {
		t.Fatalf("want 0, got %d", val)
	}

	val, _ = c.Decr(key)
	if val != -1 {
		t.Fatalf("want -1, got %d", val)
	}

	c.Set(key, 100)

	val, err = c.Incr(key)
	t.Logf("value: %d", val)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
}

func TestCtx_BitCount(t *testing.T) {
	c := AcquireCtx()
	defer c.Close()

	key := "bit_test"

	val, err := c.HasBit(key)
	t.Logf("value:%t error:%v", val, err)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	count, err := c.BitCount(key)
	t.Logf("count:%d error:%v", count, err)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	bitValue, err := c.GetBit(key, 15)
	t.Logf("bitValue:%t error:%v", bitValue, err)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	oldValue, err := c.SetBit(key, 15, true)
	t.Logf("oldValue:%t error:%v", oldValue, err)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	val, err = c.HasBit(key)
	t.Logf("value:%t error:%v", val, err)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	count, err = c.BitCount(key)
	t.Logf("count:%d error:%v", count, err)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	oldValue, err = c.SetBit(key, 15, false)
	t.Logf("oldValue:%t error:%v", oldValue, err)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	_, _ = c.SetBit(key, 0, true)
	_, _ = c.SetBit(key, 67, true)
	_, _ = c.SetBit(key, 8, true)
	_, _ = c.SetBit(key, 7, true)

	val, err = c.HasBit(key)
	t.Logf("value:%t error:%v", val, err)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	count, err = c.BitCount(key)
	t.Logf("count:%d error:%v", count, err)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}
}

func TestCtx_LList(t *testing.T) {
	c := AcquireCtx()
	defer c.Close()

	key := "list_test"

	v1Data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	v1List := make([]interface{}, len(v1Data))
	for index, value := range v1Data {
		v1List[index] = value
	}

	length, err := c.LPush(key, v1List...)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	if length != len(v1Data) {
		t.Fatalf("want 10, got %d", length)
	}

	length, err = c.LLen(key)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	if length != len(v1Data) {
		t.Fatalf("want 10, got %d", length)
	}

	valueList, err := c.LRange(key, 0, -1)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	for index, val := range valueList {
		wantValue := v1Data[len(v1Data)-index-1]
		if val.(int) != wantValue {
			t.Fatalf("want %d, got %v", wantValue, val)
		}
	}

	valueList, err = c.LRange(key, 8, -1)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	if len(valueList) != 2 {
		t.Fatalf("want %d, got %d", 2, len(valueList))
	}

	wantValue := v1Data[len(v1Data)-1-8]
	if valueList[0].(int) != v1Data[1] {
		t.Fatalf("want %d, got %v", wantValue, valueList[0])
	}

	wantValue = v1Data[len(v1Data)-1-9]
	if valueList[1].(int) != wantValue {
		t.Fatalf("want %d, got %v", wantValue, valueList[1])
	}

	value, err := c.LIndex(key, 2)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	wantValue = v1Data[len(v1Data)-1-2]
	if value.(int) != wantValue {
		t.Fatalf("want %d, got %v", wantValue, value)
	}

	value, err = c.LIndex(key, 6)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	wantValue = v1Data[len(v1Data)-1-6]
	if value.(int) != wantValue {
		t.Fatalf("want %d, got %v", wantValue, value)
	}

	err = c.LSet(key, 9, 11)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	valueList, err = c.LRange(key, -1, -1)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	if len(valueList) != 1 {
		t.Fatalf("want 1, got %d", len(valueList))
	}

	if valueList[0].(int) != 11 {
		t.Fatalf("want %d, got %v", 11, valueList[0])
	}

	err = c.LTrim(key, 0, -1)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	length, err = c.LLen(key)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	if length != len(v1Data) {
		t.Fatalf("want %d, got %d", len(v1Data), length)
	}

	err = c.LTrim(key, 1, -1)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	length, err = c.LLen(key)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	if length != len(v1Data)-1 {
		t.Fatalf("want %d, got %d", len(v1Data)-1, length)
	}

	valueList, err = c.LRange(key, 0, -1)

	val, err := c.LPop(key)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	wantValue = v1Data[len(v1Data)-2]
	if val.(int) != wantValue {
		t.Fatalf("want %d, got %v", wantValue, val)
	}

	length, err = c.LLen(key)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	if length != len(v1Data)-2 {
		t.Fatalf("want %d, got %d", len(v1Data)-2, length)
	}

	length, err = c.LPush(key, v1List...)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	if length != len(v1Data)+8 {
		t.Fatalf("want %d, got %d", len(v1Data)+8, length)
	}

	length, err = c.LLen(key)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	if length != len(v1Data)+8 {
		t.Fatalf("want %d, got %d", len(v1Data)+8, length)
	}
}

func TestCtx_LSet(t *testing.T) {
	c := AcquireCtx()
	defer c.Close()
	key := "list_set_test"

	_, _ = c.LPush(key, 7, 6, 5)
	_, _ = c.LPush(key, 4)
	_, _ = c.LPush(key, 3, 2)
	_, _ = c.LPush(key, 1)
	valueList, _ := c.LRange(key, 0, -1)
	t.Logf("value list:%+v", valueList)
	if len(valueList) != 7 {
		t.Fatalf("want 7, got %d", len(valueList))
	}

	value, _ := c.LIndex(key, 0)
	if value.(int) != 1 {
		t.Fatalf("want 1, got %+v", value)
	}

	_ = c.LSet(key, 0, value.(int)+1)
	value, _ = c.LIndex(key, 0)
	if value.(int) != 2 {
		t.Fatalf("want 2, got %+v", value)
	}

	value, _ = c.LIndex(key, 2)
	if value.(int) != 3 {
		t.Fatalf("want 3, got %+v", value)
	}

	_ = c.LSet(key, 2, value.(int)+1)
	value, _ = c.LIndex(key, 2)
	if value.(int) != 4 {
		t.Fatalf("want 4, got %+v", value)
	}

	value, _ = c.LIndex(key, 6)
	if value.(int) != 7 {
		t.Fatalf("want 7, got %+v", value)
	}

	_ = c.LSet(key, 6, value.(int)+1)
	value, _ = c.LIndex(key, 6)
	if value.(int) != 8 {
		t.Fatalf("want 8, got %+v", value)
	}

	value, _ = c.LIndex(key, -1)
	if value.(int) != 8 {
		t.Fatalf("want 8, got %+v", value)
	}

	_ = c.LSet(key, -1, value.(int)+1)
	value, _ = c.LIndex(key, -1)
	if value.(int) != 9 {
		t.Fatalf("want 9, got %+v", value)
	}

	valueList, _ = c.LRange(key, 0, -1)
	t.Logf("value list:%+v", valueList)
	if len(valueList) != 7 {
		t.Fatalf("want 7, got %d", len(valueList))
	}
}

func TestCtx_LPop(t *testing.T) {
	c := AcquireCtx()
	defer c.Close()
	key := "list_pop_test"

	length, _ := c.RPush(key, 1)
	if length != 1 {
		t.Fatalf("want 1, got %d", length)
	}

	length, _ = c.LLen(key)
	if length != 1 {
		t.Fatalf("want 1, got %d", length)
	}

	length, _ = c.LPush(key, 3, 4)
	if length != 3 {
		t.Fatalf("want 3, got %d", length)
	}

	length, _ = c.LLen(key)
	if length != 3 {
		t.Fatalf("want 3, got %d", length)
	}

	length, _ = c.RPush(key, 5, 6, 7)
	if length != 6 {
		t.Fatalf("want 6, got %d", length)
	}

	length, _ = c.LLen(key)
	if length != 6 {
		t.Fatalf("want 6, got %d", length)
	}

	item, _ := c.LPop(key)
	if item.(int) != 4 {
		t.Fatalf("want 4, got %v", item)
	}

	length, _ = c.LLen(key)
	if length != 5 {
		t.Fatalf("want 5, got %d", length)
	}

	valueList, _ := c.LRange(key, 0, -1)
	t.Logf("value list:%+v", valueList)
	if len(valueList) != 5 {
		t.Fatalf("want 5, got %d", len(valueList))
	}

	valueList, _ = c.LRange(key, 0, -2)
	t.Logf("value list:%+v", valueList)
	if len(valueList) != 4 {
		t.Fatalf("want 4, got %d", len(valueList))
	}

	valueList, _ = c.LRange(key, 0, 0)
	t.Logf("value list:%+v", valueList)
	if len(valueList) != 1 {
		t.Fatalf("want 1, got %d", len(valueList))
	}

	valueList, _ = c.LRange(key, 0, 4)
	t.Logf("value list:%+v", valueList)
	if len(valueList) != 5 {
		t.Fatalf("want 5, got %d", len(valueList))
	}

	valueList, _ = c.LRange(key, -3, -1)
	t.Logf("value list:%+v", valueList)
	if len(valueList) != 3 {
		t.Fatalf("want 3, got %d", len(valueList))
	}
}

func TestCtx_LTrim(t *testing.T) {
	c := AcquireCtx()
	defer c.Close()
	key := "list_trim_test"

	_, _ = c.LPush(key, 1, 2, 3, 4, 5, 6, 7)
	valueList, _ := c.LRange(key, 0, -1)
	t.Logf("init value list:%+v", valueList)

	err := c.LTrim(key, 0, -1)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	valueList, _ = c.LRange(key, 0, -1)
	t.Logf("trim:[0, -1] value list:%+v", valueList)
	length, _ := c.LLen(key)
	if len(valueList) != length {
		t.Fatalf("want equal, got not equal")
	}

	err = c.LTrim(key, 0, 5)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	valueList, _ = c.LRange(key, 0, -1)
	t.Logf("trim:[0, 5] value list:%+v", valueList)
	length, _ = c.LLen(key)
	if len(valueList) != length {
		t.Fatalf("want equal, got not equal")
	}

	err = c.LTrim(key, 2, 3)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	valueList, _ = c.LRange(key, 0, -1)
	t.Logf("trim:[2, 3] value list:%+v", valueList)
	length, _ = c.LLen(key)
	if len(valueList) != length {
		t.Fatalf("want equal, got not equal")
	}

	err = c.LTrim(key, 1, 3)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	valueList, _ = c.LRange(key, 0, -1)
	t.Logf("trim:[1, 3] value list:%+v", valueList)
	length, _ = c.LLen(key)
	if len(valueList) != length {
		t.Fatalf("want equal, got not equal")
	}

	err = c.LTrim(key, 1, 3)
	if err != nil {
		t.Fatalf("want nil, got %+v", err)
	}

	valueList, _ = c.LRange(key, 0, -1)
	t.Logf("trim:[1, 3] value list:%+v", valueList)
	length, _ = c.LLen(key)
	if len(valueList) != length {
		t.Fatalf("want equal, got not equal")
	}

	_, _ = c.RPush(key, 7, 6, 5, 4, 3, 2, 1)
	valueList, _ = c.LRange(key, 0, -1)
	t.Logf("value list:%+v", valueList)
	length, _ = c.LLen(key)
	if len(valueList) != length {
		t.Fatalf("want equal, got not equal")
	}
}

func TestCtx_LTrimLastest(t *testing.T) {
	c := AcquireCtx()
	defer c.Close()
	key := "list_trim_latest_test"

	var (
		num    = 100
		length int
	)

	for i := 0; i < num; i++ {
		_, err := c.RPush(key, i)
		if err != nil {
			t.Fatalf("want nil, got %+v", err)
		}

		err = c.LTrimLastest(key, 10)
		if err != nil {
			t.Fatalf("want nil, got %+v", err)
		}

		if length, _ = c.LLen(key); length > 10 {
			t.Fatalf("want less than 10, got %d", length)
		}
	}

	valueList, _ := c.LRange(key, 0, -1)
	t.Logf("value list:%+v", valueList)
}

func TestCtx_SAdd(t *testing.T) {
	c := AcquireCtx()
	defer c.Close()

	key := "set_test"

	newNum, _ := c.SAdd(key, 1, 3, 4)
	if newNum != 3 {
		t.Fatalf("want 3, got %d", newNum)
	}

	length, _ := c.SCard(key)
	if length != 3 {
		t.Fatalf("want 3, got %d", length)
	}

	newNum, _ = c.SAdd(key, 1, 3)
	if newNum != 0 {
		t.Fatalf("want 0, got %d", newNum)
	}

	length, _ = c.SCard(key)
	if length != 3 {
		t.Fatalf("want 3, got %d", length)
	}

	newNum, _ = c.SAdd(key, 5)
	if newNum != 1 {
		t.Fatalf("want 1, got %d", newNum)
	}

	length, _ = c.SCard(key)
	if length != 4 {
		t.Fatalf("want 4, got %d", length)
	}

	delNum, _ := c.SRem(key, 1)
	if delNum != 1 {
		t.Fatalf("want 1, got %d", delNum)
	}

	isMem, _ := c.SIsMember(key, 1)
	if isMem {
		t.Fatalf("want false, got:%t", isMem)
	}

	isMem, _ = c.SIsMember(key, 3)
	if !isMem {
		t.Fatalf("want true, got:%t", isMem)
	}

	length, _ = c.SCard(key)
	if length != 3 {
		t.Fatalf("want 3, got %d", length)
	}

	items, _ := c.SMembers(key)
	t.Logf("items: %+v", items)

	item, _ := c.SPop(key)
	t.Logf("pop item:%+v", item)

	item, _ = c.SRandMember(key)
	t.Logf("rand item:%+v", item)
}

func BenchmarkCtx_LPop(b *testing.B) {
	var (
		length int
		err    error
		key    = "list_bench_lpop_test"
		c      = AcquireCtx()
	)
	defer c.Close()

	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			length, err = c.LPush(key, rand.Int31())
			if err != nil {
				b.Fatalf("want nil, got %v", err)
			}
		} else {
			length, err = c.RPush(key, rand.Int31())
			if err != nil {
				b.Fatalf("want nil, got %v", err)
			}
		}

		if l, _ := c.LLen(key); l != length {
			b.Fatalf("want equal, got not equal")
		}

		if i%3 == 0 {
			_, err = c.LPop(key)
			if err != nil {
				b.Fatalf("want nil, got %v", err)
			}
		}
	}
}
