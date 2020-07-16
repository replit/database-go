package database

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	t.Parallel()
	setDBURL(t)
	const prefix = "test-json-"

	c, err := NewClient()
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
	const prefix = "test-reader-"

	c, err := NewClient()
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
