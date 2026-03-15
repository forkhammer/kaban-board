package interfaces

type PasswordServiceInterface interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) error
	GenerateSalt() (string, error)
}
