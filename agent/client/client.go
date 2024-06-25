package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
		Value: []byte(strconv.FormatInt(value, 10)),
	}

	data, _ := json.Marshal(m)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
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
		MType: "counter",
		Value: []byte(strconv.FormatFloat(value, 'f', -1, 64)),
	}

	data, _ := json.Marshal(m)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return err
}
