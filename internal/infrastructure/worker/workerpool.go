package workerpool

import (
	"context"
	"sync"

	"gostream/internal/domain"
)

type Pool struct {
	jobQueue  chan *domain.Log
	processor domain.LogProcessor
	wg        sync.WaitGroup
}

func NewPool(
	ctx context.Context,
	workerCount int,
	bufferSize int,
	processor domain.LogProcessor,
) *Pool {

	p := &Pool{
		jobQueue:  make(chan *domain.Log, bufferSize),
		processor: processor,
	}

	for i := 0; i < workerCount; i++ {
		p.wg.Add(1)
		go p.worker(ctx)
	}

	return p
}

func (p *Pool) worker(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job := <-p.jobQueue:
			_ = p.processor.Process(ctx, job)
		}
	}
}

func (p *Pool) Submit(log *domain.Log) {
	p.jobQueue <- log
}

func (p *Pool) Shutdown() {
	close(p.jobQueue)
	p.wg.Wait()
}