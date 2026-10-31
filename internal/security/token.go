package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nujabuz/music-cloud/internal/config"
	"github.com/nujabuz/music-cloud/internal/domain"
)

type tokenGenerator struct {
	cfg          config.TokenConfig
	signinMethod *jwt.SigningMethodHMAC
}

type UserClaims struct {
	jwt.RegisteredClaims
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
}

func NewTokenGenerator(cfg config.TokenConfig) domain.TokenGenerator {
	return &tokenGenerator{
		cfg:          cfg,
		signinMethod: jwt.SigningMethodHS256,
	}
}

func (tg *tokenGenerator) NewTokenUser(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"iss": tg.cfg.Issuer,
		"sub": user.Username,
		"exp": time.Now().Add(time.Second * time.Duration(tg.cfg.Expiration)).UnixMilli(),
	}

	token := jwt.NewWithClaims(tg.signinMethod, claims)
	signed, err := token.SignedString([]byte(tg.cfg.Key))

	return signed, err
}

func (tg *tokenGenerator) ParseTokenUser(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tg.cfg.Key), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return "", jwt.ErrTokenInvalidClaims
	}

	if time.Now().UnixMilli() > claims.ExpiresAt {
		return "", jwt.ErrTokenExpired
	}

	if tg.cfg.Issuer != claims.Issuer {
		return "", jwt.ErrTokenInvalidIssuer
	}

	return claims.Subject, nil
}
