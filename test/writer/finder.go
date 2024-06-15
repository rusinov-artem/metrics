package writer

import (
	"fmt"
	"strings"
	"time"
)

type Finder struct {
	Needl string
	Found bool
	Ch    chan struct{}
}

func NewFinder(substr string) *Finder {
	return &Finder{
		Needl: substr,
		Ch:    make(chan struct{}),
	}
}

func (t *Finder) Write(data []byte) (int, error) {
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

func (t *Finder) Wait(d time.Duration) error {
	select {
	case <-t.Ch:
		return nil
	case <-time.After(d):
		return fmt.Errorf("timeout. waiting for %s", t.Needl)
	}
}
