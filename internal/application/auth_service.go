package application

import (
	"context"

	"gostream/internal/domain"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo domain.UserRepository
	hasher   domain.Passwordhasher
	token    domain.Tokenmanager
}

func NewAuthService(
	userRepo domain.UserRepository,
	hasher domain.Passwordhasher,
	token domain.Tokenmanager,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		hasher:   hasher,
		token:    token,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	email string,
	password string,
	role domain.Role,
) (string, error) {

	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return "", domain.Erruseralreadyexits
	}

	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return "", err
	}

	user := &domain.User{
		ID:           uuid.New().String(),
		EMAIL:        email,
		PASSHASH: hashedPassword,
		ROLE:         string(role),
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return "", err
	}

	token, err := s.token.Generate(user.ID, domain.Role(user.ROLE))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", domain.ErrInvalidCredentials
	}

	if err := s.hasher.Compare(user.PASSHASH, password); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := s.token.Generate(user.ID, domain.Role(user.ROLE))
	if err != nil {
		return "", err
	}

	return token, nil
}