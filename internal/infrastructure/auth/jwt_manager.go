package auth

import (
	"gostream/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtmanager struct {
	secretkey     string
	tokenduretion time.Duration
}

func Newjwtmanager(secretkey string, duration time.Duration) *jwtmanager {
	return &jwtmanager{
		secretkey:     secretkey,
		tokenduretion: duration,
	}
}

type Claims struct {
	UserID string      `json:"user_id"`
	Role   domain.Role `json:"role"`
	jwt.RegisteredClaims
}

func (j *jwtmanager) Generate(userID string, role domain.Role) (string, error) {

	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenduretion)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secretkey))
}

func (j *jwtmanager) Validate(tokenStr string) (string, domain.Role, error) {

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.secretkey), nil
		},
	)
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", "", jwt.ErrTokenInvalidClaims
	}

	return claims.UserID, claims.Role, nil
}