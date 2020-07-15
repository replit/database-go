package database

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getJWT() (string, error) {
	req, err := http.NewRequest("GET", "https://database-test-jwt.kochman.repl.co", nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("test", os.Getenv("PASSWORD"))
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func TestSingleton(t *testing.T) {
	jwt, err := getJWT()
	assert.NoError(t, err)
	err = os.Setenv("REPLIT_DB_URL", jwt)
	assert.NoError(t, err)

	err = Set("test", "value")
	assert.NoError(t, err)

	val, err := Get("test")
	assert.NoError(t, err)
	assert.Equal(t, "value", val)

	err = Delete("test")
	assert.NoError(t, err)

	_, err = Get("test")
	assert.Equal(t, ErrNotFound, err)
}
