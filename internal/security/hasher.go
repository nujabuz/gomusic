package security

import (
	"github.com/nujabuz/gomusic/internal/usecase"
	"golang.org/x/crypto/bcrypt"
)

type passwordHasher struct{}

func NewPasswordHasher() usecase.PasswordHasher {
	return &passwordHasher{}
}

func (h *passwordHasher) Hash(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(hashed), err
}
