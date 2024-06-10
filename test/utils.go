package test

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

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
