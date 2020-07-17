package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ErrNotFound indicates that the requested key does not exist.
var ErrNotFound = errors.New("not found")

// client interacts with Repl.it Database. You can use it to get, set, delete,
// and list keys and their values.
type client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

// Get returns the value for the provided key. It returns ErrNotFound if the key
// does not exist.
func (c *client) Get(k string) (string, error) {
	body, err := c.GetReader(k)
	if err != nil {
		return "", err
	}
	defer body.Close()

	reader, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(reader), nil
}

// GetJSON creates or updates the provided key with the JSON serialization of
// the provided value.
func (c *client) GetJSON(k string, v interface{}) error {
	body, err := c.GetReader(k)
	if err != nil {
		return err
	}
	defer body.Close()

	return json.NewDecoder(body).Decode(&v)
}

// GetReader returns an io.ReadCloser for the value of the provided key. It
// returns ErrNotFound if the key does not exist. Callers must close the reader.
func (c *client) GetReader(k string) (io.ReadCloser, error) {
	rel := &url.URL{Path: c.baseURL.Path + "/" + k}
	u := c.baseURL.ResolveReference(rel).String()

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode <= 299 {
		return resp.Body, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, ErrNotFound
	}

	reader, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(resp.Status)
	}

	return nil, fmt.Errorf(resp.Status, string(reader))
}

// SetJSON creates or updates the provided key with the JSON serialization of
// the provided value.
func (c *client) SetJSON(k string, v interface{}) error {
	vb, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	return c.Set(k, string(vb))
}

// Set creates or updates the provided key with the provided value.
func (c *client) Set(k string, v string) error {
	rel := &url.URL{Path: c.baseURL.Path + "/" + k}
	u := c.baseURL.ResolveReference(rel)

	data := url.Values{}
	data.Set(k, v)
	ds := data.Encode()

	req, err := http.NewRequest(
		"POST",
		u.String(),
		strings.NewReader(ds),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 200 {
		return fmt.Errorf("Request failed")
	}

	return nil
}

// Delete removes the provided key from the database. No error is returned if
// the key does not exist.
func (c *client) Delete(k string) error {
	rel := &url.URL{Path: c.baseURL.Path + "/" + k}
	url := c.baseURL.ResolveReference(rel).String()
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode <= 299 {
		return nil
	}

	defer resp.Body.Close()

	reader, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(resp.Status)
	}

	return fmt.Errorf(resp.Status, string(reader))
}

// ListKeys returns a slice of all keys that begin with the provided prefix. It
// returns an empty slice if no keys match. The returned keys are sorted in
// lexicographic (string) order.
func (c *client) ListKeys(prefix string) ([]string, error) {
	rel := &url.URL{Path: c.baseURL.Path, RawQuery: "prefix=" + prefix}
	url := c.baseURL.ResolveReference(rel).String()

	var keys []string

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return keys, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return keys, err
	}

	reader, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return keys, fmt.Errorf(resp.Status)
	}

	if resp.StatusCode > 299 {
		defer resp.Body.Close()

		return keys, fmt.Errorf(resp.Status, string(reader))
	}

	strRes := string(reader)
	if len(strRes) > 0 {
		keys = strings.Split(strRes, "\n")
	}

	return keys, nil
}

// newClient returns a Client configured to use the database that is associated
// with the repl that it is running within. It does this by reading
// REPLIT_DB_URL from the environment.
func newClient() (*client, error) {
	urlStr, ok := os.LookupEnv("REPLIT_DB_URL")
	if !ok {
		return nil, fmt.Errorf("REPLIT_DB_URL not set in environment")
	}
	return newClientWithCustomURL(urlStr)
}

func newClientWithCustomURL(urlStr string) (*client, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	c := &client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: u,
	}

	return c, nil
}
