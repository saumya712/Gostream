package domain

import "context"

type LogRepository interface {
	Save(ctx context.Context, log *Log) error
}

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type MetricsStore interface {
	Increment(service string, level LogLevel)
	GetSnapshot() *Metrics
}

type LogProcessor interface {
	Process(ctx context.Context, log *Log) error
}

type Broadcaster interface {
	Broadcast(update interface{})
}

type Passwordhasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) error
}

type Tokenmanager interface {
	Generate(userid string, role Role) (string, error)
	Validate(token string) (userid string, role Role, err error)
}
