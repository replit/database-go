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

var ErrNotFound = errors.New("not found")

type Client struct {
	BaseURL    *url.URL
	httpClient *http.Client
}

func (c *Client) Get(k string) (string, error) {
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

func (c *Client) GetJSON(k string, v interface{}) error {
	body, err := c.GetReader(k)
	if err != nil {
		return err
	}
	defer body.Close()

	return json.NewDecoder(body).Decode(&v)
}

func (c *Client) GetReader(k string) (io.ReadCloser, error) {
	rel := &url.URL{Path: c.BaseURL.Path + "/" + k}
	u := c.BaseURL.ResolveReference(rel).String()

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

func (c *Client) SetJSON(k string, v interface{}) error {
	vb, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	return c.Set(k, string(vb))
}

func (c *Client) Set(k string, v string) error {
	rel := &url.URL{Path: c.BaseURL.Path + "/" + k}
	u := c.BaseURL.ResolveReference(rel)

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

func (c *Client) Delete(k string) error {
	rel := &url.URL{Path: c.BaseURL.Path + "/" + k}
	url := c.BaseURL.ResolveReference(rel).String()
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

func (c *Client) ListKeys(prefix string) ([]string, error) {
	rel := &url.URL{Path: c.BaseURL.Path, RawQuery: "prefix=" + prefix}
	url := c.BaseURL.ResolveReference(rel).String()

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

func NewClient() (*Client, error) {
	urlStr, ok := os.LookupEnv("REPLIT_DB_URL")
	if !ok {
		return nil, fmt.Errorf("REPLIT_DB_URL not set in environment")
	}
	return NewClientWithCustomUrl(urlStr)
}

func NewClientWithCustomUrl(urlStr string) (*Client, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	c := &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		BaseURL: u,
	}

	return c, nil
}
