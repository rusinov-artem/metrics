package agent

type WorkerPool struct {
	jobs chan func()
}

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

func (p *WorkerPool) Run(job func()) {
	p.jobs <- job
}

func (p *WorkerPool) Close() {
	close(p.jobs)
}
