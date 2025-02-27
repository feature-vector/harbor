package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	envFileSuffix = ".env"
)

var (
	configMap = map[string]string{}
)

func LoadEnvFromPath(configPath string) {
	err := filepath.WalkDir(configPath, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, envFileSuffix) {
			err := LoadEnvFromFile(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func Get(key string) string {
	v, exist := os.LookupEnv(key)
	if exist {
		return v
	}
	return configMap[key]
}

func LoadEnvFromFile(filepath string) error {
	file, err := os.ReadFile(filepath)
	if err != nil {
		wd, _ := os.Getwd()
		fmt.Println(fmt.Sprintf("config file [%s] not exist in working dir: %s", filepath, wd))
		return err
	}
	content := string(file)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, "=")
		if idx == -1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		configMap[key] = value
	}
	return nil
}
