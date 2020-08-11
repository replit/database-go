package database

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	t.Parallel()
	setDBURL(t)
	const prefix = "test-1-json-"

	c, err := newClient()
	assert.NoError(t, err)
	assert.NotNil(t, c)

	var str string
	err = c.GetJSON(prefix+"test", &str)
	assert.Equal(t, ErrNotFound, err)

	str = "value"
	err = c.SetJSON(prefix+"test", str)
	assert.NoError(t, err)

	str = "wrong"
	err = c.GetJSON(prefix+"test", &str)
	assert.NoError(t, err)
	assert.Equal(t, "value", str)

	err = c.Delete(prefix + "test")
	assert.NoError(t, err)
}

func TestReader(t *testing.T) {
	t.Parallel()
	setDBURL(t)
	const prefix = "test-2-reader-"

	c, err := newClient()
	assert.NoError(t, err)
	assert.NotNil(t, c)

	r, err := c.GetReader(prefix + "test")
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, r)

	err = c.Set(prefix+"test", "value")
	assert.NoError(t, err)

	r, err = c.GetReader(prefix + "test")
	assert.NoError(t, err)
	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, "value", string(b))

	err = c.Delete(prefix + "test")
	assert.NoError(t, err)
}

func TestListKeys(t *testing.T) {
	t.Parallel()
	setDBURL(t)
	const prefix = "test-3-list-keys-"

	c, err := newClient()
	assert.NoError(t, err)
	assert.NotNil(t, c)

	keys, err := c.ListKeys(prefix)
	assert.NoError(t, err)
	assert.Empty(t, keys)

	const key = prefix + "test"
	err = c.Set(key, "value")
	assert.NoError(t, err)

	keys, err = c.ListKeys(prefix)
	assert.NoError(t, err)
	assert.Equal(t, []string{key}, keys)

	err = c.Delete(key)
	assert.NoError(t, err)
}

func TestListKeysEncoding(t *testing.T) {
	t.Parallel()
	setDBURL(t)
	const prefix = "test-4-list-keys-encoding-"

	c, err := newClient()
	assert.NoError(t, err)
	assert.NotNil(t, c)

	keys, err := c.ListKeys(prefix)
	assert.NoError(t, err)
	assert.Empty(t, keys)

	const key = prefix + "\n"
	err = c.Set(key, "value")
	assert.NoError(t, err)

	keys, err = c.ListKeys(prefix)
	assert.NoError(t, err)
	assert.Equal(t, []string{key}, keys)

	err = c.Delete(key)
	assert.NoError(t, err)
}

func ExampleClient() {
	c, _ := newClient()

	c.Set("key", "value")
	val, _ := c.Get("key")
	fmt.Println(val)
	// Output: value
}
