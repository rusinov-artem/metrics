package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Destructor func()

type BufferedFileStorage struct {
	sync.Mutex
	path     string
	restore  bool
	interval int

	logger  *zap.Logger
	metrics *InMemoryMetrics
	file    *os.File
}

func NewBufferedFileStorage(logger *zap.Logger, path string, restore bool, interval int) (*BufferedFileStorage, Destructor) {
	logger = logger.With(zap.String("context", "NewBufferedFileStorage"))
	var err error
	o := &BufferedFileStorage{
		path:     path,
		restore:  restore,
		interval: interval,
		metrics:  NewInMemory(),
		logger:   logger,
	}

	o.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logger.Error("unable to open data file", zap.Error(err))
	}

	if restore {
		err = o.loadFromFile()
		if err != nil {
			logger.Error("unable to load data from file", zap.Error(err))
		}
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	go func() {
		if interval <= 0 {
			return
		}
		ticker := time.NewTicker(time.Second * time.Duration(interval))

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	destructor := func() {
		cancelFn()
		_ = o.DumpToFile()
		_ = o.file.Close()
	}

	return o, destructor
}

func (b *BufferedFileStorage) SetCounter(ctx context.Context, name string, value int64) error {
	b.Lock()
	defer b.Unlock()

	err := b.metrics.SetCounter(ctx, name, value)
	if err != nil {
		return err
	}

	if b.interval == 0 {
		return b.dumpToFile()
	}

	return nil
}

func (b *BufferedFileStorage) SetGauge(ctx context.Context, name string, value float64) error {
	b.Lock()
	defer b.Unlock()

	err := b.metrics.SetGauge(ctx, name, value)
	if err != nil {
		return err
	}

	if b.interval == 0 {
		return b.dumpToFile()
	}

	return nil
}

func (b *BufferedFileStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	b.Lock()
	defer b.Unlock()
	return b.metrics.GetCounter(ctx, name)
}

func (b *BufferedFileStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	b.Lock()
	defer b.Unlock()
	return b.metrics.GetGauge(ctx, name)
}

func (b *BufferedFileStorage) DumpToFile() error {
	b.Lock()
	defer b.Unlock()
	return b.dumpToFile()

}

func (b *BufferedFileStorage) dumpToFile() error {
	if b.file == nil {
		return fmt.Errorf("unable to dump to file, file not opened")
	}

	data, err := json.Marshal(b.metrics)
	if err != nil {
		return fmt.Errorf("unable to unmarshal data: %w", err)
	}

	err = b.file.Truncate(0)
	if err != nil {
		return fmt.Errorf("unable to truncate file: %w", err)
	}

	ret, err := b.file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("unable to go to the start of the file: %w", err)
	}

	if ret != int64(0) {
		return fmt.Errorf("offset != 0")
	}

	n, err := b.file.Write(data)
	if err != nil {
		return fmt.Errorf("unable write data to file: %w", err)
	}
	if n != len(data) {
		return fmt.Errorf("wrote %d of %d bytes: %w", n, len(data), err)
	}

	err = b.file.Sync()
	if err != nil {
		return fmt.Errorf("unable to sync file: %w", err)
	}

	return nil
}

func (b *BufferedFileStorage) LoadFromFile() error {
	b.Lock()
	defer b.Unlock()
	return b.loadFromFile()
}

func (b *BufferedFileStorage) loadFromFile() error {
	if b.file == nil {
		return fmt.Errorf("unable to load from file, file not opened")
	}

	data, err := io.ReadAll(b.file)
	if err != nil {
		return fmt.Errorf("unable to read content of the file: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	err = json.Unmarshal(data, b.metrics)
	if err != nil {
		return fmt.Errorf("unable to unmarshal data from file: %w", err)
	}

	return nil
}
