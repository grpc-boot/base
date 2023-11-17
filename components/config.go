package components

import "github.com/grpc-boot/base/core/zaplogger"

const (
	EnvPro  = `pro`
	EnvDev  = `dev`
	EnvTest = `test`
)

type Config struct {
	Name      string           `json:"name" yaml:"name"`
	Addr      string           `json:"addr" yaml:"addr"`
	PprofAddr string           `json:"pprofAddr" yaml:"pprofAddr"`
	Env       string           `json:"env" yaml:"env"`
	Ver       string           `json:"ver" yaml:"ver"`
	Logger    zaplogger.Option `json:"logger" yaml:"logger"`
	Params    JsonParam        `json:"params" yaml:"params"`
}

func (c *Config) IsEnv(env string) bool {
	return c.Env == env
}
