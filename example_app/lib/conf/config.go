package conf

import (
	"os"

	"github.com/grpc-boot/base/v2/utils"
)

var (
	config Config
)

type Config struct {
	Env string
}

// LoadConfig 通过路径加载配置文件
func LoadConfig(fileName string) error {
	confData, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	return utils.YamlDecode(utils.Bytes2String(confData), &config)
}

// GetConfig 获取配置
func GetConfig() *Config {
	return &config
}

// Env 获取环境
func Env() string {
	return config.Env
}

// IsPro 是否为生产环境
func IsPro() bool {
	return Env() == "prod"
}
