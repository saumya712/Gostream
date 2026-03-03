package processor

import (
	"context"

	"gostream/internal/domain"
)

type WorkerProcessor struct {
	metrics domain.MetricsStore
}

func NewWorkerProcessor(metrics domain.MetricsStore) *WorkerProcessor {
	return &WorkerProcessor{
		metrics: metrics,
	}
}

func (p *WorkerProcessor) Process(ctx context.Context, log *domain.Log) error {
	p.metrics.Increment(log.SERVICENAME, log.LEVEL)
	return nil
}
