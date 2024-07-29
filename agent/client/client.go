package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"

	"github.com/rusinov-artem/metrics/dto"
)

type Client struct {
	baseURL string
}

func New(baseURL string) *Client {
	baseURL = strings.TrimSuffix(baseURL, "/")
	return &Client{
		baseURL: baseURL,
	}
}

func (c *Client) SendCounter(name string, value int64) error {
	url := fmt.Sprintf("%s/update/", c.baseURL)

	m := dto.Metrics{
		ID:    name,
		MType: "counter",
		Delta: &value,
	}

	data, _ := json.Marshal(m)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return err
}

func (c *Client) SendGauge(name string, value float64) error {
	url := fmt.Sprintf("%s/update/", c.baseURL)

	m := dto.Metrics{
		ID:    name,
		MType: "gauge",
		Value: &value,
	}

	data, _ := json.Marshal(m)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	_ = resp.Body.Close()
	return err
}

func Do(fn func() error) error {
	return retry.Do(
		fn,
		retry.Attempts(3),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			if n == 1 {
				return time.Second
			}
			if n == 2 {
				return 3 * time.Second
			}
			return 5 * time.Second
		}),
	)
}
