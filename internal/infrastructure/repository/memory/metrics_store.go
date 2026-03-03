package memory

import (
	"sync"

	"gostream/internal/domain"
)

type MetricsStore struct {
	mu      sync.RWMutex
	metrics *domain.Metrics
}

func NewMetricsStore() *MetricsStore {
	return &MetricsStore{
		metrics: domain.Newmetric(),
	}
}

func (m *MetricsStore) Increment(service string, level domain.LogLevel) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.Increment(service, level)
}

func (m *MetricsStore) GetSnapshot() *domain.Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.metrics
}