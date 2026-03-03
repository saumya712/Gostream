package postgres

import (
	"context"

	"gostream/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {

	query := `
		INSERT INTO users (id, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.EMAIL,
		user.PASSHASH,
		user.ROLE,
	)

	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {

	query := `
		SELECT id, email, password_hash, role
		FROM users
		WHERE email = $1
	`

	row := r.db.QueryRow(ctx, query, email)

	var user domain.User

	err := row.Scan(
		&user.ID,
		&user.EMAIL,
		&user.PASSHASH,
		&user.ROLE,
	)

	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return &user, nil
}
