package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	req    *http.Request
	srv    *httptest.Server
	client *Client
}

func Test_Client(t *testing.T) {
	suite.Run(t, &ClientTestSuite{})
}

func (s *ClientTestSuite) SetupTest() {
	s.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.req = r
	}))

	s.client = New(fmt.Sprintf("http://%s", s.srv.Listener.Addr().String()))
}

func (s *ClientTestSuite) TearDownTest() {
	s.srv.Close()
}

func (s *ClientTestSuite) Test_CanSendCounter() {
	err := s.client.SendCounter("my_counter", 42)
	s.NoError(err)
	s.Equal("/update/counter/my_counter/42", s.req.URL.Path)
}

func (s *ClientTestSuite) Test_CanSendGauge() {
	err := s.client.SendGauge("my_gauge", 42.42)
	s.NoError(err)
	s.Equal("/update/gauge/my_gauge/42.420000", s.req.URL.Path)
}
