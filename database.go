// Package database provides methods for interacting with Repl.it Database.
// It just works if used within a repl.
package database

import (
	"sync"
)

var defaultClient = struct {
	sync.Mutex
	c *client
	o sync.Once
}{}

func getClient() (*client, error) {
	defaultClient.Lock()
	defer defaultClient.Unlock()
	if defaultClient.c == nil {
		c, err := newClient()
		if err != nil {
			return nil, err
		}
		defaultClient.c = c
	}
	return defaultClient.c, nil
}

// Get returns the value for the provided key. It returns ErrNotFound if the key
// does not exist.
func Get(key string) (string, error) {
	c, err := getClient()
	if err != nil {
		return "", err
	}

	return c.Get(key)
}

// GetJSON retrieves the JSON string data for the provided key, deserializes it
// and saves it to value pointed to by `value`.
// It returns ErrNotFound if the key does not exist.
func GetJSON(key string, value interface{}) error {
	c, err := getClient()
	if err != nil {
		return err
	}

	return c.GetJSON(key, value)
}

// Set creates or updates the provided key with the provided value.
func Set(key, value string) error {
	c, err := getClient()
	if err != nil {
		return err
	}

	return c.Set(key, value)
}

// SetJSON creates or updates the provided key with the JSON serialization of
// the provided value.
func SetJSON(key string, value interface{}) error {
	c, err := getClient()
	if err != nil {
		return err
	}

	return c.SetJSON(key, value)
}

// Delete removes the provided key.
func Delete(key string) error {
	c, err := getClient()
	if err != nil {
		return err
	}

	return c.Delete(key)
}

// ListKeys returns a slice of all keys that begin with the provided prefix.
// They are sorted in lexicographic order.
func ListKeys(prefix string) ([]string, error) {
	c, err := getClient()
	if err != nil {
		return nil, err
	}

	return c.ListKeys(prefix)
}
