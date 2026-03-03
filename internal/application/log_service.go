package application

import (
	"context"

	"gostream/internal/domain"
)

type LogDispatcher interface {
	Submit(log *domain.Log)
}

type LogService struct {
	dispatcher LogDispatcher
}

func NewLogService(dispatcher LogDispatcher) *LogService {
	return &LogService{
		dispatcher: dispatcher,
	}
}

func (s *LogService) Ingest(ctx context.Context, log *domain.Log) error {
	if !log.IsValidLevel() {
		return domain.ErrInvalidLogLevel
	}

	s.dispatcher.Submit(log)
	return nil
}
