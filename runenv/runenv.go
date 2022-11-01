package runenv

import (
	"fmt"
	"os"
	"strings"
)

// REnv custom type
type REnv string

// ToLower strings.ToLower
func (r REnv) ToLower() string {
	return strings.ToLower(string(r))
}

const (
	DefaultREnv      = Develop
	Develop     REnv = "develop"
	Test        REnv = "test"
	Gray        REnv = "gray"
	Product     REnv = "product"
)

var (
	runEnvKey = "GOWINS_RUN_ENV"
)

// Is reports whether the server is running in its env configuration
func Is(env REnv) bool {
	return GetRunEnv() == env.ToLower()
}

func Not(env REnv) bool {
	return !Is(env)
}

// IsDev reports whether the server is running in its development configuration
func IsDev() bool {
	return Is(Develop)
}

// IsTest reports whether the server is running in its testing configuration
func IsTest() bool {
	return Is(Test)
}

// IsGray reports whether the server is running in its gray configuration
func IsGray() bool {
	return Is(Gray)
}

// IsProd reports whether the server is running in its production configuration
func IsProduct() bool {
	return Is(Product)
}

// GetRunEnv the current runtime environment
func GetRunEnv() (e string) {
	if e = os.Getenv(runEnvKey); e == "" {
		// Returns a specified default value (Dev) if an empty or invalid value is detected.
		e = DefaultREnv.ToLower()
	}
	switch REnv(e) {
	case Develop, Test, Gray, Product:
	default:
		panic("unknown run environment: " + e)
	}
	return strings.ToLower(e)
}

// GetRunEnvKey the key of the runtime environment
func GetRunEnvKey() string {
	return runEnvKey
}

// SetRunEnvKey set run environment key
func SetRunEnvKey(key string) error {
	if key == "" {
		return fmt.Errorf("run environment key is empty")
	}
	runEnvKey = key
	return nil
}
