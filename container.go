package base

import (
	"sync"
)

var (
	DefaultContainer = &Container{}
)

const (
	coreConf = `core:conf`
)

type Container struct {
	container sync.Map
}

// Set 向容器中写入值
func (c *Container) Set(key string, value interface{}) {
	c.container.Store(key, value)
}

// Get 从容器中获取值
func (c *Container) Get(key string) (value interface{}, exists bool) {
	return c.container.Load(key)
}

// SetConfig 修改配置
func (c *Container) SetConfig(value *Config) {
	c.Set(coreConf, value)
}

// Config 获取配置文件
func (c *Container) Config() *Config {
	conf, ok := c.Get(coreConf)
	if ok {
		return conf.(*Config)
	}

	return nil
}
