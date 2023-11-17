package components

import (
	"math"
	"sync"
	"time"
)

const (
	abortIndex = math.MaxInt8 >> 1
)

var ctxPool = sync.Pool{
	New: func() interface{} {
		return &Context{}
	},
}

// Context 上下文
type Context struct {
	mutex sync.RWMutex

	data     map[string]interface{}
	index    int8
	handlers []func(ctx *Context)
}

// AcquireCtx 申请Context
func AcquireCtx(handlers []func(ctx *Context)) *Context {
	ctx := ctxPool.Get().(*Context)
	ctx.reset()
	ctx.handlers = handlers

	return ctx
}

// Close 释放Context
func (c *Context) Close() {
	c.reset()

	ctxPool.Put(c)
}

func (c *Context) reset() {
	c.data = nil
	c.handlers = nil
	c.index = -1
}

// Set 存储数据
func (c *Context) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.data == nil {
		c.data = make(map[string]interface{})
	}

	c.data[key] = value
}

// Get 获取数据
func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	value, exists = c.data[key]
	return
}

// GetInt64 获取int64值
func (c *Context) GetInt64(key string) (value int64, exists bool) {
	var val interface{}
	val, ok := c.Get(key)
	if !ok {
		return 0, ok
	}

	value, ok = val.(int64)

	return value, ok
}

// GetInt 获取int值
func (c *Context) GetInt(key string) (value int, exists bool) {
	var val interface{}
	val, ok := c.Get(key)
	if !ok {
		return 0, ok
	}

	value, ok = val.(int)

	return value, ok
}

// GetInt8 获取bool值
func (c *Context) GetInt8(key string) (value int8, exists bool) {
	var val interface{}
	val, ok := c.Get(key)
	if !ok {
		return 0, ok
	}

	value, ok = val.(int8)

	return value, ok
}

// GetBool 获取bool值
func (c *Context) GetBool(key string) (value bool, exists bool) {
	var val interface{}
	val, ok := c.Get(key)
	if !ok {
		return false, ok
	}

	value, ok = val.(bool)

	return value, ok
}

// GetString 获取string值
func (c *Context) GetString(key string) (value string, exists bool) {
	var val interface{}
	val, ok := c.Get(key)
	if !ok {
		return "", ok
	}

	value, ok = val.(string)

	return value, ok
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

func (c *Context) Abort() {
	c.index = abortIndex
}

/************************************/
/***** GOLANG.ORG/X/NET/CONTEXT *****/
/************************************/

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *Context) Done() <-chan struct{} {
	return nil
}

func (c *Context) Err() error {
	return nil
}

func (c *Context) Value(key interface{}) interface{} {
	if k, ok := key.(string); ok {
		v, exists := c.Get(k)
		if exists {
			return v
		}
	}

	return nil
}
