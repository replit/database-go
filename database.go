// Package database provides methods for interacting with Repl.it Database.
// It just works if used within a repl.
//
// There are two ways to use this package. The simplest is to use the top-level
// methods without creating your own client. If you want to use the more
// advanced methods, you can instantiate a client with NewClient().
package database

import (
	"sync"
)

var defaultClient = struct {
	sync.RWMutex
	c *Client
	o sync.Once
}{}

func getClient() (*Client, error) {
	if defaultClient.c == nil {
		c, err := NewClient()
		if err != nil {
			return nil, err
		}

		// TODO: we need some locking here
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

// Set creates or updates the provided key with the provided value.
func Set(key, value string) error {
	c, err := getClient()
	if err != nil {
		return err
	}

	return c.Set(key, value)
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
