package test

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	proxy *WriteProxy
	cmd   *exec.Cmd
}

func TestServer(t *testing.T) {
	suite.Run(t, &ServerTestSuite{})
}

func (this *ServerTestSuite) SetupSuite() {
	this.T().Log("SetupSuite")
}

type ProfixWriter struct {
	Prefix string
}

func (this *ProfixWriter) Write(data []byte) (int, error) {
	fmt.Printf("%s: %s", this.Prefix, string(data))
	return len(data), nil
}

type WriteProxy struct {
	W io.Writer
	sync.Mutex
}

func NewProxy() *WriteProxy {
	return &WriteProxy{
		W: &ProfixWriter{Prefix: "Empty Proxy"},
	}
}

func (this *WriteProxy) Write(data []byte) (int, error) {
	this.Lock()
	defer this.Unlock()
	return this.W.Write(data)
}

func (this *WriteProxy) SetWriter(w io.Writer) {
	this.Lock()
	defer this.Unlock()
}

func (this *WriteProxy) WaitFor(substr string) bool {
	this.Lock()
	finder := NewLookFor(substr)
	this.W = finder
	this.Unlock()
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

func (this *LookFor) Write(data []byte) (int, error) {
	fmt.Printf("%s: %s", this.Needl, string(data))
	if strings.Contains(string(data), this.Needl) {
		fmt.Printf("Found '%s' in:\n  %s\n", this.Needl, string(data))
		if !this.Found {
			this.Found = true
			close(this.Ch)
		}
	}
	return len(data), nil
}

func (this *LookFor) Wait(d time.Duration) error {
	select {
	case <-this.Ch:
		return nil
	case <-time.After(d):
		return fmt.Errorf("timeout. waiting for %s", this.Needl)
	}
}

func (this *ServerTestSuite) SetupTest() {
	var err error
	this.proxy = NewProxy()
	this.cmd = exec.Command("./srv")
	this.cmd.Stdout = this.proxy
	this.cmd.Stderr = this.proxy
	err = this.cmd.Start()
	this.NoError(err)
}

func (this *ServerTestSuite) TestMe() {
	this.AssertServerIsStarted()
	fmt.Println(this.cmd.Process.Pid)
	this.AssertStoppedCorrectly()
}

func (this *ServerTestSuite) TestLiveness() {
	this.AssertServerIsStarted()
	resp, err := http.Get("http://localhost:8080/liveness")
	this.NoError(err)
	this.Equal(200, resp.StatusCode)
	this.AssertStoppedCorrectly()
}

func (this *ServerTestSuite) AssertServerIsStarted() {
	this.T().Helper()
	this.True(this.proxy.WaitFor("Server started"), "server must start in 5 seconds")
}

func (this *ServerTestSuite) AssertStoppedCorrectly() {
	this.T().Helper()
	finder := NewLookFor("Server stopped")
	this.proxy.W = finder

	err := this.cmd.Process.Signal(syscall.SIGTERM)
	this.NoError(err)

	err = finder.Wait(5 * time.Second)
	this.NoError(err)

}
