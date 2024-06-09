package metrics

type InMemoryMetrics struct {
	Counter map[string]int64
	Guage   map[string]float64
}

func NewInMemory() *InMemoryMetrics {
	return &InMemoryMetrics{
		Counter: make(map[string]int64),
		Guage:   make(map[string]float64),
	}
}

func (t *InMemoryMetrics) SetCounter(name string, value int64) error {
	t.Counter[name] = value
	return nil
}

func (t *InMemoryMetrics) SetGuage(name string, value float64) error {
	t.Guage[name] = value
	return nil
}
