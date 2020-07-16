package database

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	setDBURL(t)

	c, err := NewClient()
	assert.NoError(t, err)
	assert.NotNil(t, c)

	var str string
	err = c.GetJSON("test", &str)
	assert.Equal(t, ErrNotFound, err)

	str = "value"
	err = c.SetJSON("test", str)
	assert.NoError(t, err)

	str = "wrong"
	err = c.GetJSON("test", &str)
	assert.NoError(t, err)
	assert.Equal(t, "value", str)

	err = c.Delete("test")
	assert.NoError(t, err)
}

func TestReader(t *testing.T) {
	setDBURL(t)

	c, err := NewClient()
	assert.NoError(t, err)
	assert.NotNil(t, c)

	r, err := c.GetReader("test")
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, r)

	err = c.Set("test", "value")
	assert.NoError(t, err)

	r, err = c.GetReader("test")
	assert.NoError(t, err)
	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, "value", string(b))

	err = c.Delete("test")
	assert.NoError(t, err)
}
