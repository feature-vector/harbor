package env

import (
	"os"
)

const (
	runtimeEnvKey = "RUNTIME_ENV"

	Production  = "prod"
	Beta        = "beta"
	Development = "dev"
	InTestCase  = "in_test_case"
)

var (
	environment = ""
)

func init() {
	environment = os.Getenv(runtimeEnvKey)
	if environment == "" {
		environment = Development
	}
}

func Environment() string {
	return environment
}

func IsProduction() bool {
	return environment == Production
}

func IsInTestCase() bool {
	return environment == InTestCase
}
