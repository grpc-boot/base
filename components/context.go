package components

import (
	"math"
	"sync"
	"time"

	"github.com/grpc-boot/base/v3/utils"
)

const (
	eventKey = `core:ent`
)

const (
	abortIndex = math.MaxInt8 >> 1
)

var ctxPool = sync.Pool{
	New: func() any {
		return &Context{}
	},
}

type Handler func(ctx *Context)

// Context 上下文
type Context struct {
	mutex sync.RWMutex

	data     map[string]any
	index    int8
	handlers []Handler
}

// AcquireCtx 申请Context
func AcquireCtx(handlers []Handler) *Context {
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
func (c *Context) Set(key string, value any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.data == nil {
		c.data = make(map[string]any)
	}

	c.data[key] = value
}

func (c *Context) SetEvent(event *Event) {
	c.Set(eventKey, event)
}

func (c *Context) Event() *Event {
	value := c.Get(eventKey, nil)
	event, _ := value.(*Event)
	return event
}

// Get 获取数据
func (c *Context) Get(key string, defaultVal any) (value any) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if _, exists := c.data[key]; !exists {
		return defaultVal
	}

	return c.data[key]
}

// GetInt 获取int64值
func (c *Context) GetInt(key string, defaultVal int64) (value int64) {
	val := c.Get(key, nil)
	if val == nil {
		return defaultVal
	}

	switch v := val.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	default:
		return defaultVal
	}
}

// GetUint 获取uint64值
func (c *Context) GetUint(key string, defaultVal uint64) (value uint64) {
	val := c.Get(key, nil)
	if val == nil {
		return defaultVal
	}

	switch v := val.(type) {
	case uint64:
		return v
	case uint:
		return uint64(v)
	case uint8:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint32:
		return uint64(v)
	default:
		return defaultVal
	}
}

// GetBool 获取bool值
func (c *Context) GetBool(key string, defaultVal bool) (value bool) {
	val := c.Get(key, nil)
	if val == nil {
		return defaultVal
	}

	switch v := val.(type) {
	case bool:
		return v
	default:
		return defaultVal
	}
}

// GetString 获取string值
func (c *Context) GetString(key, defaultVal string) (value string) {
	val := c.Get(key, nil)
	if val == nil {
		return defaultVal
	}

	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return utils.Bytes2String(v)
	default:
		return defaultVal
	}
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

func (c *Context) Value(key any) any {
	if k, ok := key.(string); ok {
		v := c.Get(k, nil)
		if v != nil {
			return v
		}
	}

	return nil
}
