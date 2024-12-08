package agent

// WorkerPool вспомогательная абстракция
// используется для отправки метрик на сервер
type WorkerPool struct {
	jobs chan func()
}

// NewWorkerPool создает WorkerPool с заданым количеством
// воркеров
func NewWorkerPool(size int) *WorkerPool {
	pool := &WorkerPool{
		jobs: make(chan func()),
	}

	for i := 0; i < size; i++ {
		go func() {
			for job := range pool.jobs {
				job()
			}
		}()
	}

	return pool
}

// Run отправляет задачу на обработку
func (p *WorkerPool) Run(job func()) {
	p.jobs <- job
}

// Close завершает обработку задач
func (p *WorkerPool) Close() {
	close(p.jobs)
}
