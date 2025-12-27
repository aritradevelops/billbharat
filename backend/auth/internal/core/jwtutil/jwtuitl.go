package jwtutil

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Dp     string `json:"dp"`
}
type claims struct {
	JwtPayload
	jwt.RegisteredClaims
}

type JwtManager struct {
	secret   string
	lifetime time.Duration
}

type AccessToken struct {
	Token    string    `json:"token"`
	Lifetime time.Time `json:"lifetime"`
}

func NewJwtManager(secret string, lifetime time.Duration) *JwtManager {
	return &JwtManager{
		secret:   secret,
		lifetime: lifetime,
	}
}

func (m *JwtManager) Sign(payload JwtPayload) (AccessToken, error) {
	claims := claims{
		JwtPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(m.secret))
	if err != nil {
		return AccessToken{}, err
	}
	return AccessToken{
		Token:    token,
		Lifetime: time.Now().Add(m.lifetime),
	}, nil
}

func (m *JwtManager) Verify(accessToken string) (*JwtPayload, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (any, error) {
		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return &claims.JwtPayload, nil
}
