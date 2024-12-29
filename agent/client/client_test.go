package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/metrics/dto"
)

type ClientTestSuite struct {
	suite.Suite
	req    *http.Request
	body   []byte
	srv    *httptest.Server
	client *Client
}

func Test_Client(t *testing.T) {
	suite.Run(t, &ClientTestSuite{})
}

func (s *ClientTestSuite) SetupTest() {
	s.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.req = r
		s.body, _ = io.ReadAll(r.Body)
	}))

	s.client = New(fmt.Sprintf("http://%s", s.srv.Listener.Addr().String()), nil)
}

func (s *ClientTestSuite) TearDownTest() {
	s.srv.Close()
}

func (s *ClientTestSuite) Test_CanSendCounter() {
	err := s.client.SendCounter("my_counter", 42)
	s.NoError(err)
	s.Equal("/update/", s.req.URL.Path)
	s.Require().NoError(err)
	m := dto.Metrics{}
	_ = json.Unmarshal(s.body, &m)
	s.Equal("my_counter", m.ID)
	s.Equal(int64(42), *m.Delta)
}

func (s *ClientTestSuite) Test_CanSendGauge() {
	err := s.client.SendGauge("my_gauge", 42.42)
	s.NoError(err)
	s.Equal("/update/", s.req.URL.Path)
	m := dto.Metrics{}
	_ = json.Unmarshal(s.body, &m)
	s.Equal("my_gauge", m.ID)
	s.InDelta(42.42, *m.Value, 0.001)
}

func (s *ClientTestSuite) Test_CanGetErrorOnSendGauge() {
	client := New(fmt.Sprintf("http://%s", "bad_url"), nil)
	err := client.SendGauge("name", 0.42)
	s.Error(err)
}

func (s *ClientTestSuite) Test_CanGetErrorOnSendCounter() {
	client := New(fmt.Sprintf("http://%s", "bad_url"), nil)
	err := client.SendCounter("name", 33)
	s.Error(err)
}
