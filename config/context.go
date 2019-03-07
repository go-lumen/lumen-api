package config

import "context"

const (
	storeKey = "config"
)

// Setter interface to set a string
type Setter interface {
	Set(string, interface{})
}

// FromContext to get value from context
func FromContext(c context.Context) *conf {
	return c.Value(storeKey).(*conf)
}

// ToContext to set value to context
func ToContext(c Setter, conf *conf) {
	c.Set(storeKey, conf)
}
