package client

import (
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	baseUrl string
}

func New(baseUrl string) *Client {
	return &Client{
		baseUrl: baseUrl,
	}
}

func (c *Client) SendCounter(name string, value int64) error {
	url := fmt.Sprintf("%s/update/counter/%s/%d", c.baseUrl, name, value)
	url = strings.TrimRight(url, "/")
	_, err := http.Post(url, "", nil)
	return err
}

func (c *Client) SendGauge(name string, value float64) error {
	url := fmt.Sprintf("%s/update/gauge/%s/%f", c.baseUrl, name, value)
	url = strings.TrimRight(url, "/")
	_, err := http.Post(url, "", nil)
	return err
}
