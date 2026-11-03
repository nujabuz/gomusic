package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nujabuz/gomusic/internal/entity"
	"github.com/nujabuz/gomusic/internal/entity/port"
)

type tokenGenerator struct {
	issuer           string
	expiresInSeconds int
	key              string
}

func NewTokenGenerator(issuer string, expiresInSeconds int, key string) port.TokenGenerator {
	return &tokenGenerator{
		issuer:           issuer,
		expiresInSeconds: expiresInSeconds,
		key:              key,
	}
}

func (tg *tokenGenerator) Generate(user *entity.User) (string, error) {
	claims := new(jwt.RegisteredClaims)
	claims.ID = fmt.Sprintf("%d", user.ID)
	claims.Subject = user.Username
	claims.Issuer = tg.issuer
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(tg.expiresInSeconds) * time.Second))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tg.key))
}

func (tg *tokenGenerator) Parse(tokenString string) (*port.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tg.key), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, jwt.ErrTokenExpired
	}

	if claims.Issuer != tg.issuer {
		return nil, jwt.ErrTokenInvalidIssuer
	}

	out := new(port.TokenClaims)
	out.Issuer = claims.Issuer
	out.Subject = claims.Subject
	out.ID = claims.ID

	return out, nil
}
