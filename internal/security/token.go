package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nujabuz/gomusic/internal/entity"
	"github.com/nujabuz/gomusic/internal/usecase"
)

type UserTokenClaims struct {
	jwt.RegisteredClaims
}

type tokenGenerator struct {
	issuer    string
	expSecond int
	key       string
}

func NewTokenGenerator(issuer string, expSeconds int, key string) usecase.TokenGenerator {
	return &tokenGenerator{
		issuer:    issuer,
		expSecond: expSeconds,
		key:       key,
	}
}

func (tg *tokenGenerator) Generate(user *entity.User) (string, error) {
	claims := new(UserTokenClaims)
	claims.ID = fmt.Sprintf("%d", user.ID)
	claims.Subject = user.Username
	claims.Issuer = tg.issuer
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(tg.expSecond) * time.Second))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tg.key))
}
