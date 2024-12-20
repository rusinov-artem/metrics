package agent

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

type Client interface {
	SendCounter(name string, value int64) error
	SendGauge(name string, value float64) error
}

// Agent в фоне отправляет метрики приложения на сервер
type Agent struct {
	sync.Mutex
	client Client

	// pollInterval интервал когда метрика обновляется
	pollInterval time.Duration

	// reporInterval интервал с которым метрика отправляется на сервер
	reportInterval time.Duration

	gauge     map[string]float64
	counter   map[string]int64
	inc       int64
	rateLimit int
}

// New создает Agent
func New(client Client, pollInterval, reportInterval time.Duration, rateLimit int) *Agent {
	defaultRateLimit := 10
	maxRateLimit := 50

	if rateLimit == 0 {
		rateLimit = defaultRateLimit
	}

	if rateLimit > maxRateLimit {
		rateLimit = maxRateLimit
	}

	return &Agent{
		client:         client,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		counter:        make(map[string]int64),
		gauge:          make(map[string]float64),
		rateLimit:      rateLimit,
	}
}

// Run запускает отправку метрик приложения на сервер
func (t *Agent) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		t.updateMetricsLoop(ctx)
	}()

	go func() {
		defer wg.Done()
		t.sendMetricsLoop(ctx)
	}()

	wg.Wait()
}

func (t *Agent) sendMetricsLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(t.reportInterval):
			t.send()
		}
	}
}

func (t *Agent) send() {
	t.Lock()
	defer t.Unlock()

	p := NewWorkerPool(t.rateLimit)
	defer p.Close()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		wg.Add(len(t.counter))
		for name, value := range t.counter {
			p.Run(func() {
				defer wg.Done()
				_ = t.client.SendCounter(name, value)
			})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		wg.Add(len(t.gauge))
		for name, value := range t.gauge {
			p.Run(func() {
				defer wg.Done()
				_ = t.client.SendGauge(name, value)
			})
		}
	}()

	wg.Wait()
}

func (t *Agent) updateMetricsLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(t.pollInterval):
			t.update()
		}
	}
}

func (t *Agent) update() {
	t.Lock()
	defer t.Unlock()

	t.inc++
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	t.gauge["Alloc"] = float64(m.Alloc)
	t.gauge["BuckHashSys"] = float64(m.BuckHashSys)
	t.gauge["Frees"] = float64(m.Frees)
	t.gauge["GCCPUFraction"] = float64(m.GCCPUFraction)
	t.gauge["GCSys"] = float64(m.GCSys)
	t.gauge["HeapAlloc"] = float64(m.HeapAlloc)
	t.gauge["HeapIdle"] = float64(m.HeapIdle)
	t.gauge["HeapInuse"] = float64(m.HeapInuse)
	t.gauge["HeapObjects"] = float64(m.HeapObjects)
	t.gauge["HeapReleased"] = float64(m.HeapReleased)
	t.gauge["HeapSys"] = float64(m.HeapSys)
	t.gauge["LastGC"] = float64(m.LastGC)
	t.gauge["Lookups"] = float64(m.Lookups)
	t.gauge["MCacheInuse"] = float64(m.MCacheInuse)
	t.gauge["MCacheSys"] = float64(m.MCacheSys)
	t.gauge["MSpanInuse"] = float64(m.MSpanInuse)
	t.gauge["MSpanSys"] = float64(m.MSpanSys)
	t.gauge["Mallocs"] = float64(m.Mallocs)
	t.gauge["NextGC"] = float64(m.NextGC)
	t.gauge["NumForcedGC"] = float64(m.NumForcedGC)
	t.gauge["NumGC"] = float64(m.NumGC)
	t.gauge["OtherSys"] = float64(m.OtherSys)
	t.gauge["PauseTotalNs"] = float64(m.PauseTotalNs)
	t.gauge["StackInuse"] = float64(m.StackInuse)
	t.gauge["StackSys"] = float64(m.StackSys)
	t.gauge["Sys"] = float64(m.Sys)
	t.gauge["TotalAlloc"] = float64(m.TotalAlloc)

	vm, err := mem.VirtualMemory()
	if err == nil {
		t.gauge["TotalMemory"] = float64(vm.Total)
		t.gauge["FreeMemory"] = float64(vm.Free)
	}

	c, err := load.Avg()
	if err == nil {
		t.gauge["CPUutilization1"] = c.Load1
	}

	t.counter["PollCount"] = t.inc + 100
	t.gauge["RandomValue"] = 42.342 * float64(t.inc)
}
