package database

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var setup sync.Once

func setDBURL(t *testing.T) {
	setup.Do(func() {
		_, ok := os.LookupEnv("RIDT_DB")
		if ok {
			fmt.Println("using file for db url")
			req, err := http.NewRequest("GET", "https://database-test-ridt-util.replit.app", nil)
			assert.NoError(t, err)

			pass, ok := os.LookupEnv("RIDT_PASSWORD")
			if !ok {
				panic("please set RIDT_PASSWORD env var")
			}
			req.SetBasicAuth("test", pass)
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			// write to /tmp/replitdb
			err = os.MkdirAll("tmp", 0755)

			err = ioutil.WriteFile(replitDBURLFile, b, 0644)

			assert.NoError(t, err)
			return
		}

		req, err := http.NewRequest("GET", "https://database-test-jwt-util.replit.app", nil)
		assert.NoError(t, err)

		pass, ok := os.LookupEnv("JWT_PASSWORD")
		if !ok {
			panic("please set JWT_PASSWORD env var")
		}
		req.SetBasicAuth("test", pass)
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = os.Setenv("REPLIT_DB_URL", string(b))
		assert.NoError(t, err)
	})
}

func TestSingleton(t *testing.T) {
	t.Parallel()
	setDBURL(t)
	const prefix = "test-singleton-"

	err := Set(prefix+"test", "value")
	assert.NoError(t, err)

	val, err := Get(prefix + "test")
	assert.NoError(t, err)
	assert.Equal(t, "value", val)

	err = Delete(prefix + "test")
	assert.NoError(t, err)

	_, err = Get(prefix + "test")
	assert.Equal(t, ErrNotFound, err)

	// listing keys
	for i := 0; i < 50; i++ {
		err = Set(fmt.Sprintf("%stest-%02d", prefix, i), "value")
		assert.NoError(t, err)
	}
	for i := 0; i < 50; i++ {
		val, err = Get(fmt.Sprintf("%stest-%02d", prefix, i))
		assert.NoError(t, err)
		assert.Equal(t, "value", val)
	}
	keys, err := ListKeys(prefix + "test")
	assert.NoError(t, err)
	assert.Len(t, keys, 50)
	for i := 0; i < 50; i++ {
		assert.Equal(t, fmt.Sprintf("%stest-%02d", prefix, i), keys[i])
	}
	for i := 0; i < 50; i++ {
		err = Delete(fmt.Sprintf("%stest-%02d", prefix, i))
		assert.NoError(t, err)
	}
}
