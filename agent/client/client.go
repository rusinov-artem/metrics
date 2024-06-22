package client

import (
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	baseURL string
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
	}
}

func (c *Client) SendCounter(name string, value int64) error {
	url := fmt.Sprintf("%s/update/counter/%s/%d", c.baseURL, name, value)
	url = strings.TrimRight(url, "/")
	resp, err := http.Post(url, "", nil)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return err
}

func (c *Client) SendGauge(name string, value float64) error {
	url := fmt.Sprintf("%s/update/gauge/%s/%f", c.baseURL, name, value)
	url = strings.TrimRight(url, "/")
	resp, err := http.Post(url, "", nil)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return err
}
