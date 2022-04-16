package config

import "context"

// StoreKey is configuration storage key
const (
	StoreKey = "config"
)

// Setter interface to set a string
type Setter interface {
	Set(string, interface{})
}

// FromContext to get value from context
func FromContext(c context.Context) *Conf {
	return c.Value(StoreKey).(*Conf)
}

// ToContext to set value to context
func ToContext(c Setter, conf *Conf) {
	c.Set(StoreKey, conf)
}
