package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"

	"github.com/rusinov-artem/metrics/dto"
)

// Client Структура которая позваляет засчитывать метрики
// Для создания этой структуры используйте функцию New
type Client struct {
	baseURL string
	Key     string
}

// New Создает структуру, которая позволяет отправлять метрики
func New(baseURL string) *Client {
	baseURL = strings.TrimSuffix(baseURL, "/")
	return &Client{
		baseURL: baseURL,
	}
}

// SendCounter отправляет значение метрики типа counter
// отправленное значение прибавляется к уже существующему
func (c *Client) SendCounter(name string, value int64) error {
	m := dto.Metrics{
		ID:    name,
		MType: "counter",
		Delta: &value,
	}

	data, _ := json.Marshal(m)

	req, err := c.newRequest(data)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return err
}

// SendGauge отправляет значение метрики типа gauge
func (c *Client) SendGauge(name string, value float64) error {
	m := dto.Metrics{
		ID:    name,
		MType: "gauge",
		Value: &value,
	}

	data, _ := json.Marshal(m)

	req, err := c.newRequest(data)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	_ = resp.Body.Close()
	return err
}

func (c *Client) newRequest(data []byte) (*http.Request, error) {
	url := fmt.Sprintf("%s/update/", c.baseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Encoding", "gzip")

	if c.Key == "" {
		return req, nil
	}

	hm := hmac.New(sha256.New, []byte(c.Key))
	hm.Write(data)
	sum := hm.Sum(nil)

	hash := base64.StdEncoding.EncodeToString(sum)
	req.Header.Set("HashSHA256", hash)

	return req, nil

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
