package storage

import (
	"context"
	"fmt"
	"sync"
)

type InMemoryMetrics struct {
	Counter map[string]int64
	Gauge   map[string]float64
	sync.Mutex
}

func NewInMemory() *InMemoryMetrics {
	return &InMemoryMetrics{
		Counter: make(map[string]int64),
		Gauge:   make(map[string]float64),
	}
}

func (t *InMemoryMetrics) SetCounter(_ context.Context, name string, value int64) error {
	t.Lock()
	defer t.Unlock()
	t.Counter[name] += value
	return nil
}

func (t *InMemoryMetrics) SetGauge(_ context.Context, name string, value float64) error {
	t.Lock()
	defer t.Unlock()
	t.Gauge[name] = value
	return nil
}

func (t *InMemoryMetrics) GetGauge(_ context.Context, name string) (float64, error) {
	v, ok := t.Gauge[name]
	if !ok {
		return 0.0, fmt.Errorf("gauge '%s' not found", name)
	}

	return v, nil
}

func (t *InMemoryMetrics) GetCounter(_ context.Context, name string) (int64, error) {
	v, ok := t.Counter[name]
	if !ok {
		return 0.0, fmt.Errorf("counter '%s' not found", name)
	}

	return v, nil
}
