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

func (this *InMemoryMetrics) SetCounter(name string, value int64) error {
	this.Counter[name] = value
	return nil
}

func (this *InMemoryMetrics) SetGuage(name string, value float64) error {
	this.Guage[name] = value
	return nil
}
