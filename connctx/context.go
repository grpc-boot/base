package connctx

import (
	"context"
	"sync"
	"time"
)

type Context interface {
	context.Context

	// Close 清理资源
	Close()

	// Del 删除key
	Del(keys ...string) (delNum int)
	// Get 获取key值
	Get(key string) (value any, exists bool)
	// Set 设置key值
	Set(key string, value any)
	// SetNx 设置key的值，当且仅当key不存在
	SetNx(key string, value any) (ok bool)
	// GetSet 设置key的值，并返回key的旧值
	GetSet(key string, value any) (old any)

	// IncrBy 将key所储存的值加上增量value
	// key存储的数据仅支持int和int64两种数据类型，否则返回类型错误
	IncrBy(key string, value int64) (newValue int64, err error)
	// DecrBy 将key所储存的值减去增量value
	// key存储的数据仅支持int和int64两种数据类型，否则返回类型错误
	DecrBy(key string, value int64) (newValue int64, err error)
	// Incr 将key所储存的值加上1
	// key存储的数据仅支持int和int64两种数据类型，否则返回类型错误
	Incr(key string) (newValue int64, err error)
	// Decr 将key所储存的值减去1
	// key存储的数据仅支持int和int64两种数据类型，否则返回类型错误
	Decr(key string) (newValue int64, err error)

	// SetBit 对key所储存的[]byte，设置或清除指定偏移量上的位(bit)
	// key存储的数据仅支持[]byte，否则返回类型错误
	SetBit(key string, offset uint16, val bool) (oldValue bool, err error)
	// GetBit 对key所储存的[]byte，获取指定偏移量上的位(bit)
	// key存储的数据仅支持[]byte，否则返回类型错误
	GetBit(key string, offset uint16) (value bool, err error)
	// HasBit 对key所储存的[]byte，获取是否含有比特位为1的位
	// key存储的数据仅支持[]byte，否则返回类型错误
	HasBit(key string) (has bool, err error)
	// BitCount 对key所储存的[]byte，获取被设置为1的比特位数量
	// key存储的数据仅支持[]byte，否则返回类型错误
	BitCount(key string) (num int, err error)

	// LPush 将一个或多个值插入到列表key的表头
	LPush(key string, items ...any) (length int, err error)
	// RPush 将一个或多个值插入到列表key的表尾
	RPush(key string, items ...any) (length int, err error)
	// LPop 移除并返回列表key的头元素
	LPop(key string) (value any, err error)
	// LIndex 返回列表key中，下标为index的元素
	LIndex(key string, index int) (value any, err error)
	// LSet 将列表key下标为index的元素的值设置为value
	LSet(key string, index int, value any) (err error)
	// LLen 返回列表key的长度
	LLen(key string) (length int, err error)
	// LRange 返回列表key中指定区间内的元素，区间以偏移量start和stop指定
	// 参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推
	LRange(key string, start, end int) (valueList []any, err error)
	// LTrim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除
	LTrim(key string, start, end int) (err error)
	// LTrimLastest 对一个列表进行修剪(trim)，保留列表key的表尾最近count条数据
	LTrimLastest(key string, count int) (err error)

	// SAdd 将一个或多个元素加入到集合key当中，已经存在于集合的元素将被忽略
	// 注意：items中不同数据类型会被认为是不同的元素
	SAdd(key string, items ...any) (newNum int, err error)
	// SCard 返回集合key中元素的数量
	SCard(key string) (total int, err error)
	// SMembers 返回集合key中的所有成员
	SMembers(key string) (items []any, err error)
	// SIsMember 判断元素是否是集合key的成员
	SIsMember(key string, item any) (isMem bool, err error)
	// SRem 移除集合key中的一个或多个元素，不存在的元素会被忽略
	SRem(key string, items ...any) (delNum int, err error)
	// SPop 移除并返回集合中的一个随机元素
	SPop(key string) (item any, err error)
	// SRandMember 返回集合中的一个随机元素
	SRandMember(key string) (item any, err error)

	// HSet 将哈希表key中的域field的值设为value
	HSet(key, field string, value any) (isCreate bool, err error)
	// HGet 返回哈希表key中给定域field的值
	HGet(key, field string) (value any, err error)
	// HGetAll 返回哈希表key中，所有的域和值
	HGetAll(key string) (value Hash, err error)
	// HDel 删除哈希表key中的一个或多个指定域，不存在的域将被忽略
	HDel(key string, fields ...string) (delNum int, err error)
	// HIncrBy 为哈希表key中的域field的值加上增量value
	// field中存储的数据仅支持int和int64两种数据类型，否则返回类型错误
	HIncrBy(key string, field string, value int64) (newValue int64, err error)
	// HLen 返回哈希表key中域的数量
	HLen(key string) (length int, err error)
	// HSetNx 将哈希表key中的域field的值设置为value，当且仅当域field不存在
	HSetNx(key, field string, value any) (ok bool, err error)
	// HSAdd 将一个或多个元素加入到哈希表key中域field集合中，已经存在于集合的元素将被忽略
	HSAdd(key, field string, items ...any) (newNum int, err error)
	// HSCard 返回哈希表key中域field集合中元素的数量
	HSCard(key, field string) (total int, err error)
	// HSMembers 返回哈希表key中域field集合中的所有成员
	HSMembers(key, field string) (items []any, err error)
	// HSIsMember 判断元素是否是哈希表key中域field集合的成员
	HSIsMember(key, field string, item any) (isMem bool, err error)
	// HSRem 移除哈希表key中域field中的一个或多个元素，不存在的元素会被忽略
	HSRem(key, field string, items ...any) (delNum int, err error)
	// HSPop 移除并返回哈希表key中域field集合中的一个随机元素
	HSPop(key, field string) (item any, err error)
	// HSRandMember 返回哈希表key中域field集合中的一个随机元素
	HSRandMember(key, field string) (item any, err error)
	// HSetBit 对哈希表key中域field所储存的[]byte，设置或清除指定偏移量上的位(bit)
	// field存储的数据仅支持[]byte，否则返回类型错误
	HSetBit(key, field string, offset uint16, val bool) (oldValue bool, err error)
	// HGetBit 对哈希表key中域field所储存的[]byte，获取指定偏移量上的位(bit)
	// field存储的数据仅支持[]byte，否则返回类型错误
	HGetBit(key, field string, offset uint16) (value bool, err error)
	// HHasBit 对哈希表key中域field所储存的[]byte，获取是否含有比特位为1的位
	// field存储的数据仅支持[]byte，否则返回类型错误
	HHasBit(key, field string) (has bool, err error)
	// HBitCount 对哈希表key中域field所储存的[]byte，获取被设置为1的比特位数量
	// field存储的数据仅支持[]byte，否则返回类型错误
	HBitCount(key, field string) (num int, err error)
}

type ctx struct {
	mutex sync.RWMutex

	data Hash
}

func newCtx() Context {
	c := &ctx{}
	return c
}

func (c *ctx) reset() {
	c.data = nil
}

func (c *ctx) setOrUpdate(handler func()) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.data == nil {
		c.data = Hash{}
	}

	handler()
}

func (c *ctx) Close() {
	c.reset()
	ctxPool.Put(c)
}

func (c *ctx) Del(keys ...string) (delNum int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.data) < 1 {
		return
	}

	return c.data.del(keys...)
}

func (c *ctx) Get(key string) (value any, exists bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.data.get(key)
}

func (c *ctx) Set(key string, value any) {
	c.setOrUpdate(func() {
		c.data.set(key, value)
	})
}

func (c *ctx) SetNx(key string, value any) (ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.data.setnx(key, value)
}

func (c *ctx) GetSet(key string, value any) (old any) {
	c.setOrUpdate(func() {
		old = c.data.getSet(key, value)
	})
	return
}

func (c *ctx) IncrBy(key string, value int64) (newValue int64, err error) {
	c.setOrUpdate(func() {
		newValue, err = c.data.incrBy(key, value)
	})
	return
}

func (c *ctx) DecrBy(key string, value int64) (newValue int64, err error) {
	return c.IncrBy(key, -value)
}

func (c *ctx) Incr(key string) (newValue int64, err error) {
	return c.IncrBy(key, 1)
}

func (c *ctx) Decr(key string) (newValue int64, err error) {
	return c.IncrBy(key, -1)
}

func (c *ctx) SetBit(key string, offset uint16, value bool) (oldValue bool, err error) {
	c.setOrUpdate(func() {
		oldValue, err = c.data.setBit(key, offset, value)
	})
	return
}

func (c *ctx) GetBit(key string, offset uint16) (value bool, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.data == nil {
		return
	}

	return c.data.getBit(key, offset)
}

func (c *ctx) HasBit(key string) (has bool, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.data == nil {
		return
	}

	return c.data.hasBit(key)
}

func (c *ctx) BitCount(key string) (num int, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.data == nil {
		return
	}

	return c.data.bitCount(key)
}

func (c *ctx) LPush(key string, items ...any) (length int, err error) {
	if len(items) < 1 {
		return
	}

	c.setOrUpdate(func() {
		value, exists := c.data[key]
		if !exists {
			l := &list{}
			l.prepend(items...)
			length = l.length
			c.data[key] = l
			return
		}

		val, ok := value.(*list)
		if !ok {
			err = ErrType
			return
		}

		val.prepend(items...)
		length = val.length
	})

	return
}

func (c *ctx) RPush(key string, items ...any) (length int, err error) {
	if len(items) < 1 {
		return
	}

	c.setOrUpdate(func() {
		value, exists := c.data[key]
		if !exists {
			l := &list{}
			l.append(items...)
			length = l.length
			c.data[key] = l
			return
		}

		val, ok := value.(*list)
		if !ok {
			err = ErrType
			return
		}

		val.append(items...)
		length = val.length
	})

	return
}

func (c *ctx) LPop(key string) (value any, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.data) < 1 {
		return
	}

	val, exists := c.data[key]
	if !exists {
		return
	}

	v, ok := val.(*list)
	if !ok {
		err = ErrType
		return
	}

	value = v.lpop()
	return
}

func (c *ctx) LIndex(key string, index int) (value any, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.data) < 1 {
		return
	}

	val, exists := c.data[key]
	if !exists {
		return
	}

	v, ok := val.(*list)
	if !ok {
		err = ErrType
		return
	}

	value = v.index(index)
	return
}

func (c *ctx) LSet(key string, index int, value any) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.data) < 1 {
		return
	}

	val, exists := c.data[key]
	if !exists {
		return
	}

	v, ok := val.(*list)
	if !ok {
		err = ErrType
		return
	}

	return v.set(index, value)
}

func (c *ctx) LLen(key string) (length int, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.data) < 1 {
		return
	}

	val, exists := c.data[key]
	if !exists {
		return
	}

	v, ok := val.(*list)
	if !ok {
		err = ErrType
		return
	}

	length = v.length
	return
}

func (c *ctx) LRange(key string, start, end int) (valueList []any, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.data) < 1 {
		return
	}

	val, exists := c.data[key]
	if !exists {
		return
	}

	v, ok := val.(*list)
	if !ok {
		err = ErrType
		return
	}

	return v.lrange(start, end)
}

func (c *ctx) LTrim(key string, start, end int) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.data) < 1 {
		return
	}

	val, exists := c.data[key]
	if !exists {
		return
	}

	v, ok := val.(*list)
	if !ok {
		err = ErrType
		return
	}

	v.trim(start, end)
	return
}

func (c *ctx) LTrimLastest(key string, count int) (err error) {
	return c.LTrim(key, -count, -1)
}

func (c *ctx) SAdd(key string, items ...any) (newNum int, err error) {
	c.setOrUpdate(func() {
		newNum, err = c.data.sAdd(key, items...)
	})

	return
}

func (c *ctx) SCard(key string) (total int, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.data == nil {
		return
	}

	return c.data.sCard(key)
}

func (c *ctx) SMembers(key string) (items []any, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.data == nil {
		return
	}

	return c.data.sMembers(key)
}

func (c *ctx) SIsMember(key string, item any) (isMem bool, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.data == nil {
		return
	}

	return c.data.sIsMember(key, item)
}

func (c *ctx) SRem(key string, items ...any) (delNum int, err error) {
	c.setOrUpdate(func() {
		delNum, err = c.data.sRem(key, items...)
	})

	return
}

func (c *ctx) SPop(key string) (item any, err error) {
	c.setOrUpdate(func() {
		item, err = c.data.sPop(key)
	})

	return
}

func (c *ctx) SRandMember(key string) (item any, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.data == nil {
		return
	}

	return c.data.sRandMember(key)
}

func (c *ctx) HSet(key, field string, value any) (isCreate bool, err error) {
	c.setOrUpdate(func() {
		val, exists := c.data[key]
		if !exists {
			isCreate = true

			c.data[key] = Hash{
				field: value,
			}

			return
		}

		if _, ok := val.(Hash); !ok {
			err = ErrType
			return
		}

		isCreate = c.data[key].(Hash).set(field, value)
	})
	return
}

func (c *ctx) HMSet(key string, kv Hash) (err error) {
	if len(kv) < 1 {
		return
	}

	c.setOrUpdate(func() {
		value, exists := c.data[key]
		if !exists {
			c.data[key] = kv
			return
		}

		if _, ok := value.(Hash); !ok {
			err = ErrType
			return
		}

		c.data[key].(Hash).mset(kv)
	})
	return
}

func (c *ctx) HGet(key, field string) (value any, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.data == nil {
		return
	}

	val, exists := c.data[key]
	if !exists {
		return
	}

	v, ok := val.(Hash)
	if !ok {
		err = ErrType
		return
	}

	value, _ = v.get(field)
	return
}

func (c *ctx) HGetAll(key string) (value Hash, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.data) < 1 {
		return
	}

	val, exists := c.data[key]
	if !exists {
		return
	}

	if _, ok := val.(Hash); !ok {
		err = ErrType
		return
	}

	value = c.data[key].(Hash).getall()
	return
}

func (c *ctx) HDel(key string, fields ...string) (delNum int, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.data) < 1 {
		return
	}

	value, exists := c.data[key]
	if !exists {
		return
	}

	if _, ok := value.(Hash); !ok {
		err = ErrType
		return
	}

	delNum = c.data[key].(Hash).del(fields...)
	return
}

func (c *ctx) HIncrBy(key string, field string, value int64) (newValue int64, err error) {
	c.setOrUpdate(func() {
		val, exists := c.data[key]
		if !exists {
			c.data[key] = Hash{
				field: value,
			}

			newValue = value
			return
		}

		if _, ok := val.(Hash); !ok {
			err = ErrType
			return
		}

		newValue, err = c.data[key].(Hash).incrBy(field, value)
	})

	return
}

func (c *ctx) HLen(key string) (length int, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.data) < 1 {
		return
	}

	value, exists := c.data[key]
	if !exists {
		return
	}

	val, ok := value.(Hash)
	if !ok {
		err = ErrType
		return
	}

	length = val.len()
	return
}

func (c *ctx) HSetNx(key, field string, value any) (ok bool, err error) {
	c.setOrUpdate(func() {
		val, exists := c.data[key]
		if !exists {
			ok = true
			c.data[key] = Hash{
				field: value,
			}
			return
		}

		if _, yes := val.(Hash); !yes {
			err = ErrType
			return
		}

		ok = c.data[key].(Hash).setnx(field, value)
	})

	return
}

func (c *ctx) HSAdd(key, field string, items ...any) (newNum int, err error) {
	if len(items) < 1 {
		return
	}

	c.setOrUpdate(func() {
		value, exists := c.data[key]
		if !exists {
			newNum = len(items)
			s := make(set, len(items))
			s.add(items...)

			c.data[key] = Hash{
				field: s,
			}

			return
		}

		if _, ok := value.(Hash); !ok {
			err = ErrType
			return
		}

		newNum, err = c.data[key].(Hash).sAdd(field, items...)
	})

	return
}

func (c *ctx) HSCard(key, field string) (total int, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	value, exists := c.data[key]
	if !exists {
		return
	}

	val, ok := value.(Hash)
	if !ok {
		err = ErrType
		return
	}

	return val.sCard(field)
}

func (c *ctx) HSMembers(key, field string) (items []any, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	value, exists := c.data[key]
	if !exists {
		return
	}

	val, ok := value.(Hash)
	if !ok {
		err = ErrType
		return
	}

	return val.sMembers(field)
}

func (c *ctx) HSIsMember(key, field string, item any) (isMem bool, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	value, exists := c.data[key]
	if !exists {
		return
	}

	val, ok := value.(Hash)
	if !ok {
		err = ErrType
		return
	}

	return val.sIsMember(field, item)
}

func (c *ctx) HSRem(key, field string, items ...any) (delNum int, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.data) < 1 {
		return
	}

	value, exists := c.data[key]
	if !exists {
		return
	}

	if _, ok := value.(Hash); !ok {
		err = ErrType
		return
	}

	return c.data[key].(Hash).sRem(field, items...)
}

func (c *ctx) HSPop(key, field string) (item any, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.data) < 1 {
		return
	}

	value, exists := c.data[key]
	if !exists {
		return
	}

	if _, ok := value.(Hash); !ok {
		err = ErrType
		return
	}

	return c.data[key].(Hash).sPop(field)
}

func (c *ctx) HSRandMember(key, field string) (item any, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	value, exists := c.data[key]
	if !exists {
		return
	}

	val, ok := value.(Hash)
	if !ok {
		err = ErrType
		return
	}

	return val.sRandMember(field)
}

func (c *ctx) HSetBit(key, field string, offset uint16, val bool) (oldValue bool, err error) {
	c.setOrUpdate(func() {
		value, exists := c.data[key]
		if !exists {
			if !val {
				return
			}

			h := Hash{}
			oldValue, _ = h.setBit(field, offset, val)
			c.data[key] = h
			return
		}

		if _, ok := value.(Hash); !ok {
			err = ErrType
			return
		}

		oldValue, err = c.data[key].(Hash).setBit(field, offset, val)
	})
	return
}

func (c *ctx) HGetBit(key, field string, offset uint16) (value bool, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.data) < 1 {
		return
	}

	val, exists := c.data[key]
	if !exists {
		return
	}

	v, ok := val.(Hash)
	if !ok {
		err = ErrType
		return
	}

	return v.getBit(field, offset)
}

func (c *ctx) HHasBit(key, field string) (has bool, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.data) < 1 {
		return
	}

	val, exists := c.data[key]
	if !exists {
		return
	}

	v, ok := val.(Hash)
	if !ok {
		err = ErrType
		return
	}

	return v.hasBit(field)
}

func (c *ctx) HBitCount(key, field string) (num int, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.data) < 1 {
		return
	}

	val, exists := c.data[key]
	if !exists {
		return
	}

	v, ok := val.(Hash)
	if !ok {
		err = ErrType
		return
	}

	return v.bitCount(field)
}

/************************************/
/***** GOLANG.ORG/X/NET/CONTEXT *****/
/************************************/

func (c *ctx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *ctx) Done() <-chan struct{} {
	return nil
}

func (c *ctx) Err() error {
	return nil
}

func (c *ctx) Value(key any) any {
	if k, ok := key.(string); ok {
		v, exists := c.Get(k)
		if exists {
			return v
		}
	}

	return nil
}
