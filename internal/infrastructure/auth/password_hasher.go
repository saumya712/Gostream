package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type Bcrypthasher struct {
	cost int
}

func Newbcrypthasher(cost int) *Bcrypthasher {
	return &Bcrypthasher{
		cost: cost,
	}
}

func (b *Bcrypthasher) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (b *Bcrypthasher) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}
