package config

import (
	"context"

	"github.com/spf13/viper"
)

type conf struct {
	*viper.Viper
}

// New allows to create a viper configuration
func New(viper *viper.Viper) *conf {
	return &conf{viper}
}

// GetString allows to retrieve a specific string
func GetString(c context.Context, key string) string {
	return FromContext(c).GetString(key)
}

// GetBool allows to retrieve a specific bool
func GetBool(c context.Context, key string) bool {
	return FromContext(c).GetBool(key)
}

// GetInt allows to retrieve a specific int
func GetInt(c context.Context, key string) int {
	return FromContext(c).GetInt(key)
}

// Set allows to set a value
func Set(c context.Context, key string, value interface{}) {
	FromContext(c).Set(key, value)
}
