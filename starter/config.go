package starter

import "github.com/feature-vector/harbor/base/conf"

func InitConfig(configPath string) {
	conf.LoadEnvFromPath(configPath)
}
