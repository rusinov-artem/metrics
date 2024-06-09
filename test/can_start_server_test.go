package test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/rusinov-artem/metrics/agent/client"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	proxy *WriteProxy
	cmd   *exec.Cmd
}

func TestServer(t *testing.T) {
	BinTest := os.Getenv("BIN_TEST")
	if BinTest != "TRUE" {
		t.Skip("enable this test by env BIN_TEST=TRUE")
	}
	suite.Run(t, &ServerTestSuite{})
}

func (t *ServerTestSuite) SetupSuite() {
	t.T().Log("SetupSuite")
}

type ProfixWriter struct {
	Prefix string
}

func (t *ProfixWriter) Write(data []byte) (int, error) {
	fmt.Printf("%s: %s", t.Prefix, string(data))
	return len(data), nil
}

type WriteProxy struct {
	w io.Writer
	sync.Mutex
}

func NewProxy() *WriteProxy {
	return &WriteProxy{
		w: &ProfixWriter{Prefix: "Empty Proxy"},
	}
}

func (t *WriteProxy) Write(data []byte) (int, error) {
	t.Lock()
	defer t.Unlock()
	return t.w.Write(data)
}

func (t *WriteProxy) SetWriter(w io.Writer) {
	t.Lock()
	defer t.Unlock()
	t.w = w
}

func (t *WriteProxy) WaitFor(substr string) bool {
	finder := NewLookFor(substr)
	t.SetWriter(finder)
	err := finder.Wait(5 * time.Second)
	return err == nil
}

type LookFor struct {
	Needl string
	Found bool
	Ch    chan struct{}
}

func NewLookFor(substr string) *LookFor {
	return &LookFor{
		Needl: substr,
		Ch:    make(chan struct{}),
	}
}

func (t *LookFor) Write(data []byte) (int, error) {
	fmt.Printf("%s: %s", t.Needl, string(data))
	if strings.Contains(string(data), t.Needl) {
		fmt.Printf("Found '%s' in:\n  %s\n", t.Needl, string(data))
		if !t.Found {
			t.Found = true
			close(t.Ch)
		}
	}
	return len(data), nil
}

func (t *LookFor) Wait(d time.Duration) error {
	select {
	case <-t.Ch:
		return nil
	case <-time.After(d):
		return fmt.Errorf("timeout. waiting for %s", t.Needl)
	}
}

func (t *ServerTestSuite) SetupTest() {
	var err error
	t.proxy = NewProxy()
	t.cmd = exec.Command("./srv")
	t.cmd.Stdout = t.proxy
	t.cmd.Stderr = t.proxy
	err = t.cmd.Start()
	t.NoError(err)
	t.AssertServerIsStarted()
}

func (t *ServerTestSuite) TearDownTest() {
	t.AssertStoppedCorrectly()
}

func (t *ServerTestSuite) TestLiveness() {
	resp, err := http.Get("http://localhost:8080/liveness")
	defer func() { _ = resp.Body.Close() }()
	t.NoError(err)
	t.Equal(200, resp.StatusCode)
}

func (t *ServerTestSuite) TestCanSetCounterByClient() {
	c := client.New("http://localhost:8080")

	finder := NewLookFor("my_counter")
	t.proxy.SetWriter(finder)

	err := c.SendCounter("my_counter", 42)
	t.NoError(err)

	err = finder.Wait(time.Second)
	t.NoError(err)
}

func (t *ServerTestSuite) TestCanSetCounterByGauge() {
	c := client.New("http://localhost:8080")

	finder := NewLookFor("my_gauge")
	t.proxy.SetWriter(finder)

	err := c.SendGauge("my_gauge", 42.42)
	t.NoError(err)

	err = finder.Wait(time.Second)
	t.NoError(err)
}

func (t *ServerTestSuite) AssertServerIsStarted() {
	t.T().Helper()
	t.True(t.proxy.WaitFor("Server started"), "server must start in 5 seconds")
}

func (t *ServerTestSuite) AssertStoppedCorrectly() {
	t.T().Helper()
	finder := NewLookFor("Server stopped")
	t.proxy.SetWriter(finder)

	err := t.cmd.Process.Signal(syscall.SIGTERM)
	t.NoError(err)

	err = finder.Wait(5 * time.Second)
	t.NoError(err)

}
