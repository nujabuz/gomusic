package security

import (
	"errors"

	"github.com/nujabuz/music-cloud/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type passwordHasher struct{}

func NewPasswordHasher() domain.PasswordHasher {
	return &passwordHasher{}
}

func (passwordHasher) Hash(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(hashed), err
}
func (passwordHasher) Compare(plain string, hashed string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return domain.ErrInvalidUser
	}
	return err
}
