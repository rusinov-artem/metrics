package test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/metrics/agent/client"
	"github.com/rusinov-artem/metrics/test/writer"
)

type ServerTestSuite struct {
	suite.Suite
	proxy         *writer.WriterProxy
	cmd           *exec.Cmd
	serverAddress string
}

func TestServer(t *testing.T) {
	BinTest := os.Getenv("BIN_TEST")
	if BinTest != "TRUE" {
		t.Skip("enable this test by env BIN_TEST=TRUE")
	}

	t.Run("server without arguments", func(t *testing.T) {
		suite.Run(t, &ServerTestSuite{})
	})

	t.Run("server with custom address", func(t *testing.T) {
		suite.Run(t, &ServerTestSuite{serverAddress: "127.0.0.1:9999"})
	})
}

func (t *ServerTestSuite) SetupSuite() {
	t.T().Log("SetupSuite")
}

func (t *ServerTestSuite) SetupTest() {
	var err error
	t.proxy = writer.NewProxy()
	cmdName, cmdArgs := t.buildCmd()
	t.cmd = exec.Command(cmdName, cmdArgs...)
	t.cmd.Stdout = t.proxy
	t.cmd.Stderr = t.proxy
	err = t.cmd.Start()
	t.NoError(err)
	t.AssertServerIsStarted()

	if t.serverAddress == "" {
		t.serverAddress = "0.0.0.0:8080"
	}
}

func (t *ServerTestSuite) TearDownTest() {
	t.AssertStoppedCorrectly()
}

func (t *ServerTestSuite) TestLiveness() {
	resp, err := http.Get(t.livenessURL())
	t.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	t.NoError(err)
	t.Equal(200, resp.StatusCode)
}

func (t *ServerTestSuite) TestCanSetCounterByClient() {
	c := client.New(t.baseURL())

	finder := writer.NewFinder("my_counter")
	t.proxy.SetWriter(finder)

	err := c.SendCounter("my_counter", 42)
	t.NoError(err)

	err = finder.Wait(time.Second)
	t.NoError(err)
}

func (t *ServerTestSuite) TestCanGetCounterValue() {
	c := client.New(t.baseURL())

	err := c.SendCounter("my_counter", 42)
	t.NoError(err)

	resp, err := http.Get(t.baseURL() + "/value/counter/my_counter")
	t.NoError(err)
	b, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	t.Equal("42", string(b))
}

func (t *ServerTestSuite) TestCanSetCounterByGauge() {
	c := client.New(t.baseURL())

	finder := writer.NewFinder("my_gauge")
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
	finder := writer.NewFinder("Server stopped")
	t.proxy.SetWriter(finder)

	err := t.cmd.Process.Signal(syscall.SIGTERM)
	t.NoError(err)

	err = finder.Wait(5 * time.Second)
	t.NoError(err)
}

func (t *ServerTestSuite) buildCmd() (string, []string) {
	cmd := []string{"./srv"}

	if t.serverAddress != "" {
		cmd = append(cmd, "-a", t.serverAddress)
	}

	return cmd[0], cmd[1:]
}

func (t *ServerTestSuite) livenessURL() string {
	return fmt.Sprintf("http://%s/liveness", t.serverAddress)
}

func (t *ServerTestSuite) baseURL() string {
	return fmt.Sprintf("http://%s", t.serverAddress)
}
