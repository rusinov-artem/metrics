package writer

import (
	"io"
	"sync"
	"time"
)

type WriterProxy struct {
	w io.Writer
	sync.Mutex
}

func NewProxy() *WriterProxy {
	return &WriterProxy{
		w: &PrefixWriter{Prefix: "Empty Proxy"},
	}
}

func (t *WriterProxy) Write(data []byte) (int, error) {
	t.Lock()
	defer t.Unlock()
	return t.w.Write(data)
}

func (t *WriterProxy) SetWriter(w io.Writer) {
	t.Lock()
	defer t.Unlock()
	t.w = w
}

func (t *WriterProxy) WaitFor(substr string) bool {
	finder := NewFinder(substr)
	t.SetWriter(finder)
	err := finder.Wait(5 * time.Second)
	return err == nil
}
