package base

import "github.com/grpc-boot/base/core/zaplogger"

const (
	EnvPro  = `pro`
	EnvDev  = `dev`
	EnvTest = `test`
)

type Config struct {
	Name   string           `json:"name" yaml:"name"`
	Addr   string           `json:"addr" yaml:"addr"`
	Env    string           `json:"env" yaml:"env"`
	Ver    string           `json:"ver" yaml:"ver"`
	Logger zaplogger.Option `json:"logger" yaml:"logger"`
}

func (c *Config) IsPro() bool {
	return c.Env == EnvPro
}

func (c *Config) IsTest() bool {
	return c.Env == EnvTest
}

func (c *Config) IsDev() bool {
	return c.Env == EnvDev
}
