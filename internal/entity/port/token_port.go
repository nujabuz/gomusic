package port

import "github.com/nujabuz/gomusic/internal/entity"

// TokenClaims output from TokenGenerator.Parse
type TokenClaims struct {
	Issuer  string
	Subject string
	ID      string
}

// TokenGenerator port
type TokenGenerator interface {
	Generate(user *entity.User) (string, error)
	Parse(tokenString string) (*TokenClaims, error)
}
