package agent

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type AgentTestSuite struct {
	suite.Suite
	client *FakeClient
}

func TestAgent(t *testing.T) {
	suite.Run(t, &AgentTestSuite{})
}

func (s *AgentTestSuite) SetupTest() {
	s.client = NewFakeClient()
}

func (s *AgentTestSuite) Test_CreateAgent() {
	agent := New(s.client, 10*time.Millisecond, 100*time.Millisecond)
	ctx, closeFn := context.WithTimeout(context.Background(), 2*time.Second)
	defer closeFn()
	agent.Run(ctx)
	s.Greater(s.client.sendGuageExecuted, 10)
	s.Greater(s.client.sendCounterExecuted, 10)
}

type FakeClient struct {
	sendCounterExecuted int
	sendGuageExecuted   int
}

func (f *FakeClient) SendCounter(name string, value int64) error {
	f.sendCounterExecuted++
	return nil
}
func (f *FakeClient) SendGauge(name string, value float64) error {
	f.sendGuageExecuted++
	return nil
}

func NewFakeClient() *FakeClient {
	return &FakeClient{}
}
