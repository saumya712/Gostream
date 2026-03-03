package domain

import "time"

type LogLevel string

const (
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

type Log struct {
	ID          string
	SERVICENAME string
	LEVEL       LogLevel
	MESSAGE     string
	TIMESTAMP   time.Time
}

func (l *Log) IsValidLevel() bool {
	switch l.LEVEL {
	case LevelInfo, LevelWarn, LevelError:
		return true
	default:
		return false
	}
}
