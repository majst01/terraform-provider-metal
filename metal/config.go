package metal

import (
	metalgo "github.com/metal-pod/metal-go"
)

type Config struct {
	hmac string
}

// Client returns a new client for accessing Metal's API.
func (c *Config) Client() (*metalgo.Driver, error) {
	return metalgo.NewDriver("metal.test.fi-ts.io", "", c.hmac)
}
