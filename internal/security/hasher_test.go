package security_test

import (
	"strings"
	"testing"

	"github.com/nujabuz/music-cloud/internal/domain"
	"github.com/nujabuz/music-cloud/internal/security"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasherTestSuite struct {
	suite.Suite
	hasher domain.PasswordHasher
}

// SetupTest initializes the hasher before each test
func (s *PasswordHasherTestSuite) SetupTest() {
	s.hasher = security.NewPasswordHasher()
}

// TestHash_Success tests generating a valid bcrypt hash
func (s *PasswordHasherTestSuite) TestHash_Success() {
	plain := "mySecretPassword123"

	hashed, err := s.hasher.Hash(plain)

	s.NoError(err)
	s.NotEmpty(hashed)
	s.NotEqual(plain, hashed)
	s.NoError(bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)))
}

// TestHash_PasswordTooLong tests error handling when password exceeds bcrypt length limit
func (s *PasswordHasherTestSuite) TestHash_PasswordTooLong() {
	// bcrypt has a maximum input length limit of 72 bytes
	longPassword := strings.Repeat("a", 73)

	hashed, err := s.hasher.Hash(longPassword)

	s.Error(err)
	s.Empty(hashed)
}

// TestCompare_Success tests successful matching of plain text and hashed password
func (s *PasswordHasherTestSuite) TestCompare_Success() {
	plain := "securePass"
	hashed, err := s.hasher.Hash(plain)
	s.Require().NoError(err)

	err = s.hasher.Compare(plain, hashed)

	s.NoError(err)
}

// TestCompare_Mismatch tests error translation when passwords do not match
func (s *PasswordHasherTestSuite) TestCompare_Mismatch() {
	plain := "correctPassword"
	wrong := "wrongPassword"
	hashed, err := s.hasher.Hash(plain)
	s.Require().NoError(err)

	err = s.hasher.Compare(wrong, hashed)

	s.ErrorIs(err, domain.ErrInvalidUser)
}

// TestCompare_InvalidHash tests error handling with an invalid bcrypt hash format
func (s *PasswordHasherTestSuite) TestCompare_InvalidHash() {
	err := s.hasher.Compare("password", "invalid-hash-string")

	s.Error(err)
	s.NotErrorIs(err, domain.ErrInvalidUser)
}

func TestPasswordHasherTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordHasherTestSuite))
}
