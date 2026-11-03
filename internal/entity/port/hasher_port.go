package port

// PasswordHasher port
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hashedPassword string, password string) error
}
