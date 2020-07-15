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

func Get(key string) (string, error) {
	c, err := getClient()
	if err != nil {
		return "", err
	}

	return c.Get(key)
}

func Set(key, value string) error {
	c, err := getClient()
	if err != nil {
		return err
	}

	return c.Set(key, value)
}

func Delete(key string) error {
	c, err := getClient()
	if err != nil {
		return err
	}

	return c.Delete(key)
}

func ListKeys(prefix string) ([]string, error) {
	c, err := getClient()
	if err != nil {
		return nil, err
	}

	return c.ListKeys(prefix)
}
